// Copyright 2023 The go-ethereum Authors
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

package eth

// PIDAPI provides an API to control the PID controller for advanced fee management.
type PIDAPI struct {
	e *Ethereum
}

// NewPIDAPI creates a new PIDAPI instance.
func NewPIDAPI(e *Ethereum) *PIDAPI {
	return &PIDAPI{e}
}

// SetPIDEnabled enables or disables the PID controller
func (api *PIDAPI) SetPIDEnabled(enabled bool) bool {
	api.e.Miner().SetPIDEnabled(enabled)
	return true
}

// UpdatePIDParameters allows runtime adjustment of PID parameters
func (api *PIDAPI) UpdatePIDParameters(kp, ki, kd, maxChange float64) bool {
	api.e.Miner().UpdatePIDParameters(kp, ki, kd, maxChange)
	return true
}

// UpdatePIDExternalPressure receives pressure signals from batcher via RPC
func (api *PIDAPI) UpdatePIDExternalPressure(daPressure, l1Influence float64) bool {
	api.e.Miner().UpdatePIDExternalPressure(daPressure, l1Influence)
	return true
}

// GetPIDStatus returns current PID controller status
func (api *PIDAPI) GetPIDStatus() map[string]interface{} {
	return api.e.Miner().GetPIDStatus()
}
