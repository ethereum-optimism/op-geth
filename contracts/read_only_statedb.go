package contracts

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

// ReadOnlyStateDB wraps a StateDB to prevent modifications. Using it without copying it is safe.
//
// Only the following methods will be passed through to the original StateDB:
//
//	GetState(common.Address, common.Hash) common.Hash
//	GetCodeHash(common.Address) common.Hash
//	GetCode(common.Address) []byte
//
// Gas calculations based on ReadOnlyStateDB will be wrong because the accessed storage slots and addresses are not tracked.
type ReadOnlyStateDB struct {
	vm.StateDB
}

func (r *ReadOnlyStateDB) CreateAccount(common.Address) {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) SubBalance(common.Address, *uint256.Int) {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) AddBalance(_ common.Address, amount *uint256.Int) {
	if amount.Cmp(new(uint256.Int)) == 0 {
		// Adding zero is safe, so we can return here
		return
	}
	panic("not implemented")
}

func (r *ReadOnlyStateDB) GetBalance(common.Address) *uint256.Int {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) GetNonce(common.Address) uint64 {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) SetNonce(common.Address, uint64) {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) SetCode(common.Address, []byte) {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) GetCodeSize(common.Address) int {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) AddRefund(uint64) {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) SubRefund(uint64) {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) GetRefund() uint64 {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) GetCommittedState(common.Address, common.Hash) common.Hash {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) SetState(common.Address, common.Hash, common.Hash) {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) GetTransientState(addr common.Address, key common.Hash) common.Hash {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) SetTransientState(addr common.Address, key, value common.Hash) {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) SelfDestruct(common.Address) {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) HasSelfDestructed(common.Address) bool {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) Selfdestruct6780(common.Address) {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) Exist(common.Address) bool {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) Empty(common.Address) bool {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) AddressInAccessList(addr common.Address) bool {
	// We don't track access lists
	return false
}

func (r *ReadOnlyStateDB) SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool) {
	// We don't track access lists
	return false, false
}

func (r *ReadOnlyStateDB) AddAddressToAccessList(addr common.Address) {
}

func (r *ReadOnlyStateDB) AddSlotToAccessList(addr common.Address, slot common.Hash) {
}

func (r *ReadOnlyStateDB) Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) RevertToSnapshot(int) {
	// No changes can be done, so reverting is a noop.
}

func (r *ReadOnlyStateDB) Snapshot() int {
	// We use id 0 for the single state this immutable StateDB can have.
	return 0
}

func (r *ReadOnlyStateDB) AddLog(*types.Log) {
	panic("not implemented")
}

func (r *ReadOnlyStateDB) AddPreimage(common.Hash, []byte) {
	panic("not implemented")
}
