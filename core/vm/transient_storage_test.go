package vm

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestRevert(t *testing.T) {

	ts := NewTransientStorage()

	address1 := common.HexToAddress("0xa")

	slot1 := common.HexToHash("0x0")
	value1 := common.HexToHash("0x1")

	ts.Set(address1, slot1, value1)

	address2 := common.HexToAddress("0xa")

	slot2 := common.HexToHash("0x1")
	value2 := common.HexToHash("0x2")

	ts.Set(address2, slot2, value2)

	ts.CheckPoint()

	err := ts.Revert()

	if err != nil {
		t.Error(err)
	}

	if ts.current[address1][slot1] != value1 {
		t.Error("First transient storage entry could not be found")
	}

	if ts.current[address2][slot2] != value2 {
		t.Error("Second transient storage entry could not be found")
	}

}

func TestRevert2(t *testing.T) {

	ts := NewTransientStorage()

	address1 := common.HexToAddress("0xa")

	slot1 := common.HexToHash("0x0")
	value1 := common.HexToHash("0x1")

	ts.Set(address1, slot1, value1)

	address2 := common.HexToAddress("0xa")

	slot2 := common.HexToHash("0x1")
	value2 := common.HexToHash("0x2")

	ts.Set(address2, slot2, value2)

	ts.CheckPoint()

	err := ts.Revert()

	if err != nil {
		t.Error(err)
	}

	if ts.current[address1][slot1] != value1 {
		t.Error("First transient storage entry could not be found")
	}

	if ts.current[address2][slot2] != value2 {
		t.Error("Second transient storage entry could not be found")
	}

}
