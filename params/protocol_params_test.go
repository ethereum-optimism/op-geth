// Copyright 2026 The go-ethereum Authors
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

package params

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestEIP8282SystemContracts(t *testing.T) {
	t.Parallel()
	// Addresses and runtime-code digests from go-ethereum v1.17.7:
	// https://github.com/ethereum/go-ethereum/blob/v1.17.7/params/protocol_params.go
	for _, tc := range []struct {
		name        string
		address     common.Address
		code        []byte
		wantAddress string
		wantSHA256  string
	}{
		{
			name:        "builder deposit",
			address:     BuilderDepositAddress,
			code:        BuilderDepositCode,
			wantAddress: "0x0000BFF46984E3725691FA540A8C7589300D8282",
			wantSHA256:  "2c49dcf745b1304f3dac0ea7487eae6d8fd07812ada980d542f79e8e5e53eb8d",
		},
		{
			name:        "builder exit",
			address:     BuilderExitAddress,
			code:        BuilderExitCode,
			wantAddress: "0x000064D678505AD48F8CCB093BC65613800E8282",
			wantSHA256:  "c889ed88730d157d192aae28c2dee61324d0df3bd01ff0078386808b4adb27aa",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if want := common.HexToAddress(tc.wantAddress); tc.address != want {
				t.Errorf("address = %s, want %s", tc.address, want)
			}
			if got := fmt.Sprintf("%x", sha256.Sum256(tc.code)); got != tc.wantSHA256 {
				t.Errorf("runtime-code SHA-256 = %s, want %s", got, tc.wantSHA256)
			}
		})
	}
}
