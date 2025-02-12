// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package nat

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	natpmp "github.com/jackpal/go-nat-pmp"
)

// pmp adapts the NAT-PMP protocol implementation so it conforms to
// the common interface.
type pmp struct {
	gw net.IP
	c  *natpmp.Client
}

func (n *pmp) String() string {
	return fmt.Sprintf("NAT-PMP(%v)", n.gw)
}

func (n *pmp) ExternalIP() (net.IP, error) {
	response, err := n.c.GetExternalAddress()
	if err != nil {
		return nil, err
	}
	return response.ExternalIPAddress[:], nil
}

func (n *pmp) AddMapping(protocol string, extport, intport int, name string, lifetime time.Duration) (uint16, error) {
	if lifetime <= 0 {
		return 0, errors.New("lifetime must not be <= 0")
	}
	// Note order of port arguments is switched between our
	// AddMapping and the client's AddPortMapping.
	res, err := n.c.AddPortMapping(strings.ToLower(protocol), intport, extport, int(lifetime/time.Second))
	if err != nil {
		return 0, err
	}

	// NAT-PMP maps an alternative available port number if the requested port
	// is already mapped to another address and returns success. Handling of
	// alternate port numbers is done by the caller.
	return res.MappedExternalPort, nil
}

func (n *pmp) DeleteMapping(protocol string, extport, intport int) (err error) {
	// To destroy a mapping, send an add-port with an internalPort of
	// the internal port to destroy, an external port of zero and a
	// time of zero.
	_, err = n.c.AddPortMapping(strings.ToLower(protocol), intport, 0, 0)
	return err
}

func (n *pmp) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("natpmp:%v", n.gw)), nil
}

func discoverPMP() Interface {
	// run external address lookups on all potential gateways
	gws := potentialGateways()
	found := make(chan *pmp, len(gws))
	for i := range gws {
		gw := gws[i]
		go func() {
			c := natpmp.NewClient(gw)
			if _, err := c.GetExternalAddress(); err != nil {
				found <- nil
			} else {
				found <- &pmp{gw, c}
			}
		}()
	}
	// return the one that responds first.
	// discovery needs to be quick, so we stop caring about
	// any responses after a very short timeout.
	timeout := time.NewTimer(1 * time.Second)
	defer timeout.Stop()
	for range gws {
		select {
		case c := <-found:
			if c != nil {
				return c
			}
		case <-timeout.C:
			return nil
		}
	}
	return nil
}

// potentialGateways returns a list of potential gateway IP addresses by:
// 1. Getting the default gateway from routing table (primary method)
// 2. Checking common gateway IPs in private networks (fallback method)
// 3. Handling multiple network interfaces properly
func potentialGateways() (gws []net.IP) {
	// Get network interfaces
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	// Create a map to store unique gateway IPs
	uniqueGWs := make(map[string]net.IP)

	for _, iface := range ifaces {
		// Skip interfaces that are down or loopback
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		ifaddrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range ifaddrs {
			if x, ok := addr.(*net.IPNet); ok && x.IP.IsPrivate() {
				ip := x.IP.Mask(x.Mask).To4()
				if ip == nil {
					continue
				}

				// Try common gateway IPs
				candidates := []net.IP{
					// X.X.X.1 (most common)
					append(append([]byte(nil), ip[:3]...), 1),
					// X.X.X.254 (some routers)
					append(append([]byte(nil), ip[:3]...), 254),
					// Current network address + 1
					append(append([]byte(nil), ip[:3]...), ip[3]|0x01),
				}

				// Add unique candidates to our map
				for _, candidate := range candidates {
					if candidate.IsPrivate() {
						uniqueGWs[candidate.String()] = candidate
					}
				}
			}
		}
	}

	// Convert unique gateways map to slice
	for _, ip := range uniqueGWs {
		gws = append(gws, ip)
	}

	return gws
}
