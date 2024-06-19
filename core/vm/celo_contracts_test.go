package vm

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/addresses"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
	"golang.org/x/crypto/sha3"
)

func makeTestHeaderHash(number *big.Int) common.Hash {
	preimage := append([]byte("fakeheader"), common.LeftPadBytes(number.Bytes()[:], 32)...)
	return common.Hash(sha3.Sum256(preimage))
}

func makeTestHeader(number *big.Int) *types.Header {
	return &types.Header{
		ParentHash: makeTestHeaderHash(new(big.Int).Sub(number, common.Big1)),
		Number:     number,
		GasUsed:    params.DefaultGasLimit / 2,
		Time:       number.Uint64() * 5,
	}
}

var testHeader = makeTestHeader(big.NewInt(10000))

var vmBlockCtx = BlockContext{
	CanTransfer: func(db StateDB, addr common.Address, amount *uint256.Int) bool {
		return db.GetBalance(addr).Cmp(amount) >= 0
	},
	Transfer: func(db StateDB, a1, a2 common.Address, i *uint256.Int) {
		panic("transfer: not implemented")
	},
	GetHash: func(u uint64) common.Hash {
		panic("getHash: not implemented")
	},
	Coinbase:    common.Address{},
	BlockNumber: new(big.Int).Set(testHeader.Number),
	Time:        testHeader.Time,
}

var vmTxCtx = TxContext{
	GasPrice: common.Big1,
	Origin:   common.HexToAddress("a11ce"),
}

// Create a global mock EVM for use in the following tests.
var mockEVM = &EVM{
	chainConfig: params.TestChainConfig,
	Context:     vmBlockCtx,
	TxContext:   vmTxCtx,
}

var mockPrecompileContext = NewContext(common.HexToAddress("1337"), mockEVM)

func TestPrecompileTransfer(t *testing.T) {
	type args struct {
		input []byte
		ctx   *celoPrecompileContext
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr string
	}{
		{
			name: "Test transfer with invalid caller",
			args: args{
				input: []byte(""),
				ctx:   mockPrecompileContext,
			},
			wantErr:     true,
			expectedErr: "unable to call transfer from unpermissioned address",
		}, {
			name: "Test transfer with short input",
			args: args{
				input: []byte("0000"),
				ctx:   NewContext(addresses.CeloTokenAddress, mockEVM),
			},
			wantErr:     true,
			expectedErr: "invalid input length",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &transfer{}
			_, err := c.Run(tt.args.input, tt.args.ctx)
			if tt.wantErr {
				if err == nil {
					t.Errorf("transfer.Run() expected error = %v", tt.expectedErr)
				} else if err.Error() != tt.expectedErr {
					t.Errorf("transfer.Run() error = %v, expected %v", err, tt.wantErr)
				}
			}
		})
	}
}
