package overlay

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/holiman/uint256"
)

type OverlayCreateTracer struct {
	contractAddress common.Address
	creator         common.Address
	isCapturing     bool
	foundCreator    bool
	code            []byte
	gasCap          uint64
	err             error
	resultCode      []byte
	evm             *vm.EVM
}

// Transaction level
func (ct *OverlayCreateTracer) CaptureTxStart(gasLimit uint64) {
}
func (ct *OverlayCreateTracer) CaptureTxEnd(restGas uint64) {
}

// Top call frame
func (ct *OverlayCreateTracer) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	ct.evm = env
	if create && to == ct.contractAddress {
		log.Debug("[CaptureStart]", "from", from, "to", to)
		ct.foundCreator = true
		ct.creator = from
	}
}
func (ct *OverlayCreateTracer) CaptureEnd(output []byte, gasUsed uint64, err error) {
}

// Rest of call frames
func (ct *OverlayCreateTracer) CaptureEnter(typ vm.OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
	if ct.isCapturing {
		return
	}

	if (typ == vm.CREATE || typ == vm.CREATE2) && to == ct.contractAddress {
		log.Debug("[CaptureEnter]", "to", to, "typ", typ)
		ct.foundCreator = true
		ct.creator = from
		if ct.code != nil {
			ct.isCapturing = true
			_, _, _, err := ct.evm.OverlayCreate(vm.AccountRef(from), vm.NewCodeAndHash(ct.code), ct.gasCap, uint256.MustFromBig(value), to, typ, true /* incrementNonce */)
			if err != nil {
				ct.err = err
			} else {
				ct.resultCode = ct.evm.StateDB.GetCode(ct.contractAddress)
			}
		}
	}
}
func (ct *OverlayCreateTracer) CaptureExit(output []byte, gasUsed uint64, err error) {
}

// Opcode level
func (ct *OverlayCreateTracer) CaptureState(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
}
func (ct *OverlayCreateTracer) CaptureFault(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
}
