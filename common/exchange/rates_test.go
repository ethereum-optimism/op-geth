package exchange

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

var (
	currA         = common.HexToAddress("0xA")
	currB         = common.HexToAddress("0xB")
	currX         = common.HexToAddress("0xF")
	exchangeRates = common.ExchangeRates{
		currA: big.NewRat(47, 100),
		currB: big.NewRat(45, 100),
	}
)

func TestCompareFees(t *testing.T) {
	type args struct {
		val1         *big.Int
		feeCurrency1 *common.Address
		val2         *big.Int
		feeCurrency2 *common.Address
	}
	tests := []struct {
		name       string
		args       args
		wantResult int
		wantErr    bool
	}{
		// Native currency
		{
			name: "Same amount of native currency",
			args: args{
				val1:         big.NewInt(1),
				feeCurrency1: nil,
				val2:         big.NewInt(1),
				feeCurrency2: nil,
			},
			wantResult: 0,
		}, {
			name: "Different amounts of native currency 1",
			args: args{
				val1:         big.NewInt(2),
				feeCurrency1: nil,
				val2:         big.NewInt(1),
				feeCurrency2: nil,
			},
			wantResult: 1,
		}, {
			name: "Different amounts of native currency 2",
			args: args{
				val1:         big.NewInt(1),
				feeCurrency1: nil,
				val2:         big.NewInt(5),
				feeCurrency2: nil,
			},
			wantResult: -1,
		},
		// Mixed currency
		{
			name: "Same amount of mixed currency",
			args: args{
				val1:         big.NewInt(1),
				feeCurrency1: nil,
				val2:         big.NewInt(1),
				feeCurrency2: &currA,
			},
			wantResult: -1,
		}, {
			name: "Different amounts of mixed currency 1",
			args: args{
				val1:         big.NewInt(100),
				feeCurrency1: nil,
				val2:         big.NewInt(47),
				feeCurrency2: &currA,
			},
			wantResult: 0,
		}, {
			name: "Different amounts of mixed currency 2",
			args: args{
				val1:         big.NewInt(45),
				feeCurrency1: &currB,
				val2:         big.NewInt(100),
				feeCurrency2: nil,
			},
			wantResult: 0,
		},
		// Two fee currencies
		{
			name: "Same amount of same currency",
			args: args{
				val1:         big.NewInt(1),
				feeCurrency1: &currA,
				val2:         big.NewInt(1),
				feeCurrency2: &currA,
			},
			wantResult: 0,
		}, {
			name: "Different amounts of same currency 1",
			args: args{
				val1:         big.NewInt(3),
				feeCurrency1: &currA,
				val2:         big.NewInt(1),
				feeCurrency2: &currA,
			},
			wantResult: 1,
		}, {
			name: "Different amounts of same currency 2",
			args: args{
				val1:         big.NewInt(1),
				feeCurrency1: &currA,
				val2:         big.NewInt(7),
				feeCurrency2: &currA,
			},
			wantResult: -1,
		}, {
			name: "Different amounts of different currencies 1",
			args: args{
				val1:         big.NewInt(47),
				feeCurrency1: &currA,
				val2:         big.NewInt(45),
				feeCurrency2: &currB,
			},
			wantResult: 0,
		}, {
			name: "Different amounts of different currencies 2",
			args: args{
				val1:         big.NewInt(48),
				feeCurrency1: &currA,
				val2:         big.NewInt(45),
				feeCurrency2: &currB,
			},
			wantResult: 1,
		}, {
			name: "Different amounts of different currencies 3",
			args: args{
				val1:         big.NewInt(47),
				feeCurrency1: &currA,
				val2:         big.NewInt(46),
				feeCurrency2: &currB,
			},
			wantResult: -1,
		},
		// Unregistered fee currency
		{
			name: "Different amounts of different currencies",
			args: args{
				val1:         big.NewInt(1),
				feeCurrency1: &currA,
				val2:         big.NewInt(1),
				feeCurrency2: &currX,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompareValue(exchangeRates, tt.args.val1, tt.args.feeCurrency1, tt.args.val2, tt.args.feeCurrency2)

			if tt.wantErr && err == nil {
				t.Error("Expected error in CompareValue()")
			}
			if got != tt.wantResult {
				t.Errorf("CompareValue() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}
