// Copyright 2024 The go-ethereum Authors
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

package version

import (
	"regexp"
	"strconv"
)

// Upstream geth version
const (
	Major = 1        // Major version component of the current release
	Minor = 16       // Minor version component of the current release
	Patch = 4        // Patch version component of the current release
	Meta  = "stable" // Version metadata to append to the version string
)

// OPGeth is the version of op-geth
var (
	OPGethMajor = 0          // Major version component of the current release
	OPGethMinor = 1          // Minor version component of the current release
	OPGethPatch = 0          // Patch version component of the current release
	OPGethMeta  = "untagged" // Version metadata to append to the version string
)

// This is set at build-time by the linker when the build is done by build/ci.go.
var gitTag string

// Override the version variables if the gitTag was set at build time.
var _ = func() (_ string) {
	semver := regexp.MustCompile(`^v([0-9]+)\.([0-9]+)\.([0-9]+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+[0-9A-Za-z-]+)?$`)
	version := semver.FindStringSubmatch(gitTag)
	if version == nil {
		return
	}
	if version[4] == "" {
		version[4] = "stable"
	}
	OPGethMajor, _ = strconv.Atoi(version[1])
	OPGethMinor, _ = strconv.Atoi(version[2])
	OPGethPatch, _ = strconv.Atoi(version[3])
	OPGethMeta = version[4]
	return
}()
