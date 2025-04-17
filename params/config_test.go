// Copyright 2017 The go-ethereum Authors
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
	"fmt"
	"math"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCheckCompatible(t *testing.T) {
	type test struct {
		stored, new   *ChainConfig
		headBlock     uint64
		headTimestamp uint64
		wantErr       *ConfigCompatError

		genesisTimestamp *uint64
	}
	tests := []test{
		{stored: AllEthashProtocolChanges, new: AllEthashProtocolChanges, headBlock: 0, headTimestamp: 0, wantErr: nil},
		{stored: AllEthashProtocolChanges, new: AllEthashProtocolChanges, headBlock: 0, headTimestamp: uint64(time.Now().Unix()), wantErr: nil},
		{stored: AllEthashProtocolChanges, new: AllEthashProtocolChanges, headBlock: 100, wantErr: nil},
		{
			stored:    &ChainConfig{EIP150Block: big.NewInt(10)},
			new:       &ChainConfig{EIP150Block: big.NewInt(20)},
			headBlock: 9,
			wantErr:   nil,
		},
		{
			stored:    AllEthashProtocolChanges,
			new:       &ChainConfig{HomesteadBlock: nil},
			headBlock: 3,
			wantErr: &ConfigCompatError{
				What:          "Homestead fork block",
				StoredBlock:   big.NewInt(0),
				NewBlock:      nil,
				RewindToBlock: 0,
			},
		},
		{
			stored:    AllEthashProtocolChanges,
			new:       &ChainConfig{HomesteadBlock: big.NewInt(1)},
			headBlock: 3,
			wantErr: &ConfigCompatError{
				What:          "Homestead fork block",
				StoredBlock:   big.NewInt(0),
				NewBlock:      big.NewInt(1),
				RewindToBlock: 0,
			},
		},
		{
			stored:    &ChainConfig{HomesteadBlock: big.NewInt(30), EIP150Block: big.NewInt(10)},
			new:       &ChainConfig{HomesteadBlock: big.NewInt(25), EIP150Block: big.NewInt(20)},
			headBlock: 25,
			wantErr: &ConfigCompatError{
				What:          "EIP150 fork block",
				StoredBlock:   big.NewInt(10),
				NewBlock:      big.NewInt(20),
				RewindToBlock: 9,
			},
		},
		{
			stored:    &ChainConfig{ConstantinopleBlock: big.NewInt(30)},
			new:       &ChainConfig{ConstantinopleBlock: big.NewInt(30), PetersburgBlock: big.NewInt(30)},
			headBlock: 40,
			wantErr:   nil,
		},
		{
			stored:    &ChainConfig{ConstantinopleBlock: big.NewInt(30)},
			new:       &ChainConfig{ConstantinopleBlock: big.NewInt(30), PetersburgBlock: big.NewInt(31)},
			headBlock: 40,
			wantErr: &ConfigCompatError{
				What:          "Petersburg fork block",
				StoredBlock:   nil,
				NewBlock:      big.NewInt(31),
				RewindToBlock: 30,
			},
		},
		{
			stored:        &ChainConfig{ShanghaiTime: newUint64(10)},
			new:           &ChainConfig{ShanghaiTime: newUint64(20)},
			headTimestamp: 9,
			wantErr:       nil,
		},
		{
			stored:        &ChainConfig{ShanghaiTime: newUint64(10)},
			new:           &ChainConfig{ShanghaiTime: newUint64(20)},
			headTimestamp: 25,
			wantErr: &ConfigCompatError{
				What:         "Shanghai fork timestamp",
				StoredTime:   newUint64(10),
				NewTime:      newUint64(20),
				RewindToTime: 9,
			},
		},
		{
			stored:           &ChainConfig{CanyonTime: newUint64(10)},
			new:              &ChainConfig{CanyonTime: newUint64(20)},
			headTimestamp:    25,
			genesisTimestamp: newUint64(2),
			wantErr: &ConfigCompatError{
				What:         "Canyon fork timestamp",
				StoredTime:   newUint64(10),
				NewTime:      newUint64(20),
				RewindToTime: 9,
			},
		},
		{
			stored:           &ChainConfig{CanyonTime: newUint64(10)},
			new:              &ChainConfig{CanyonTime: newUint64(20)},
			headTimestamp:    25,
			genesisTimestamp: nil,
			wantErr: &ConfigCompatError{
				What:         "Canyon fork timestamp",
				StoredTime:   newUint64(10),
				NewTime:      newUint64(20),
				RewindToTime: 9,
			},
		},
		{
			stored:           &ChainConfig{CanyonTime: newUint64(10)},
			new:              &ChainConfig{CanyonTime: newUint64(20)},
			headTimestamp:    25,
			genesisTimestamp: newUint64(24),
			wantErr:          nil,
		},
		{
			stored:           &ChainConfig{HoloceneTime: newUint64(10)},
			new:              &ChainConfig{HoloceneTime: newUint64(20)},
			headTimestamp:    25,
			genesisTimestamp: newUint64(15),
			wantErr: &ConfigCompatError{
				What:         "Holocene fork timestamp",
				StoredTime:   newUint64(10),
				NewTime:      newUint64(20),
				RewindToTime: 9,
			},
		},
		{
			stored:           &ChainConfig{HoloceneTime: newUint64(10)},
			new:              &ChainConfig{HoloceneTime: newUint64(20)},
			headTimestamp:    15,
			genesisTimestamp: newUint64(5),
			wantErr: &ConfigCompatError{
				What:         "Holocene fork timestamp",
				StoredTime:   newUint64(10),
				NewTime:      newUint64(20),
				RewindToTime: 9,
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := test.stored.CheckCompatible(test.new, test.headBlock, test.headTimestamp, test.genesisTimestamp)
			if !reflect.DeepEqual(err, test.wantErr) {
				t.Errorf("error mismatch:\nstored: %v\nnew: %v\nheadBlock: %v\nheadTimestamp: %v\nerr: %v\nwant: %v", test.stored, test.new, test.headBlock, test.headTimestamp, err, test.wantErr)
			}
		})
	}
}

func TestConfigRules(t *testing.T) {
	c := &ChainConfig{
		LondonBlock:  new(big.Int),
		ShanghaiTime: newUint64(500),
	}
	var stamp uint64
	if r := c.Rules(big.NewInt(0), true, stamp); r.IsShanghai {
		t.Errorf("expected %v to not be shanghai", stamp)
	}
	stamp = 500
	if r := c.Rules(big.NewInt(0), true, stamp); !r.IsShanghai {
		t.Errorf("expected %v to be shanghai", stamp)
	}
	stamp = math.MaxInt64
	if r := c.Rules(big.NewInt(0), true, stamp); !r.IsShanghai {
		t.Errorf("expected %v to be shanghai", stamp)
	}
}

func TestTimestampCompatError(t *testing.T) {
	require.Equal(t, new(ConfigCompatError).Error(), "")

	errWhat := "Shanghai fork timestamp"
	require.Equal(t, newTimestampCompatError(errWhat, nil, newUint64(1681338455)).Error(),
		"mismatching Shanghai fork timestamp in database (have timestamp nil, want timestamp 1681338455, rewindto timestamp 1681338454)")

	require.Equal(t, newTimestampCompatError(errWhat, newUint64(1681338455), nil).Error(),
		"mismatching Shanghai fork timestamp in database (have timestamp 1681338455, want timestamp nil, rewindto timestamp 1681338454)")

	require.Equal(t, newTimestampCompatError(errWhat, newUint64(1681338455), newUint64(600624000)).Error(),
		"mismatching Shanghai fork timestamp in database (have timestamp 1681338455, want timestamp 600624000, rewindto timestamp 600623999)")

	require.Equal(t, newTimestampCompatError(errWhat, newUint64(0), newUint64(1681338455)).Error(),
		"mismatching Shanghai fork timestamp in database (have timestamp 0, want timestamp 1681338455, rewindto timestamp 0)")
}

func TestConfigRulesRegolith(t *testing.T) {
	c := &ChainConfig{
		RegolithTime: newUint64(500),
		LondonBlock:  new(big.Int),
		Optimism:     &OptimismConfig{},
	}
	var stamp uint64
	if r := c.Rules(big.NewInt(0), true, stamp); r.IsOptimismRegolith {
		t.Errorf("expected %v to not be regolith", stamp)
	}
	stamp = 500
	if r := c.Rules(big.NewInt(0), true, stamp); !r.IsOptimismRegolith {
		t.Errorf("expected %v to be regolith", stamp)
	}
	stamp = math.MaxInt64
	if r := c.Rules(big.NewInt(0), true, stamp); !r.IsOptimismRegolith {
		t.Errorf("expected %v to be regolith", stamp)
	}
}

func TestCheckOptimismValidity(t *testing.T) {
	validOpConfig := &OptimismConfig{
		EIP1559Denominator:       10,
		EIP1559Elasticity:        50,
		EIP1559DenominatorCanyon: newUint64(250),
	}

	tests := []struct {
		name    string
		config  *ChainConfig
		wantErr *string
	}{
		{
			name: "valid",
			config: &ChainConfig{
				Optimism:     validOpConfig,
				CanyonTime:   newUint64(100),
				ShanghaiTime: newUint64(100),
				CancunTime:   newUint64(200),
				EcotoneTime:  newUint64(200),
				PragueTime:   newUint64(300),
				IsthmusTime:  newUint64(300),
			},
			wantErr: nil,
		},
		{
			name: "zero EIP1559Denominator",
			config: &ChainConfig{
				Optimism: &OptimismConfig{
					EIP1559Denominator: 0,
					EIP1559Elasticity:  50,
				},
			},
			wantErr: ptr("zero EIP1559Denominator"),
		},
		{
			name: "zero EIP1559Elasticity",
			config: &ChainConfig{
				Optimism: &OptimismConfig{
					EIP1559Denominator: 10,
					EIP1559Elasticity:  0,
				},
			},
			wantErr: ptr("zero EIP1559Elasticity"),
		},
		{
			name: "missing EIP1559DenominatorCanyon",
			config: &ChainConfig{
				Optimism: &OptimismConfig{
					EIP1559Denominator: 10,
					EIP1559Elasticity:  50,
				},
				CanyonTime: newUint64(100),
			},
			wantErr: ptr("missing or zero EIP1559DenominatorCanyon"),
		},
		{
			name: "ShanghaiTime not equal to CanyonTime",
			config: &ChainConfig{
				Optimism:     validOpConfig,
				ShanghaiTime: newUint64(100),
				CanyonTime:   newUint64(200),
			},
			wantErr: ptr("ShanghaiTime (100) must equal CanyonTime (200)"),
		},
		{
			name: "CancunTime not equal to EcotoneTime",
			config: &ChainConfig{
				Optimism:    validOpConfig,
				CancunTime:  newUint64(200),
				EcotoneTime: newUint64(300),
			},
			wantErr: ptr("CancunTime (200) must equal EcotoneTime (300)"),
		},
		{
			name: "PragueTime not equal to IsthmusTime",
			config: &ChainConfig{
				Optimism:    validOpConfig,
				PragueTime:  newUint64(300),
				IsthmusTime: newUint64(400),
			},
			wantErr: ptr("PragueTime (300) must equal IsthmusTime (400)"),
		},
		{
			name: "nil ShanghaiTime",
			config: &ChainConfig{
				Optimism:   validOpConfig,
				CanyonTime: newUint64(200),
			},
			wantErr: ptr("ShanghaiTime (<nil>) must equal CanyonTime (200)"),
		},
		{
			name: "nil CancunTime",
			config: &ChainConfig{
				Optimism:    validOpConfig,
				EcotoneTime: newUint64(300),
			},
			wantErr: ptr("CancunTime (<nil>) must equal EcotoneTime (300)"),
		},
		{
			name: "nil PragueTime",
			config: &ChainConfig{
				Optimism: &OptimismConfig{
					EIP1559Denominator:       10,
					EIP1559Elasticity:        50,
					EIP1559DenominatorCanyon: newUint64(250),
				},
				IsthmusTime: newUint64(400),
			},
			wantErr: ptr("PragueTime (<nil>) must equal IsthmusTime (400)"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.CheckOptimismValidity()
			if tt.wantErr != nil {
				require.EqualError(t, err, *tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func ptr[T any](t T) *T {
	return &t
}
