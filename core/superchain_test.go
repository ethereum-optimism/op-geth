package core

import (
	"testing"

	"github.com/ethereum-optimism/superchain-registry/superchain"
)

func TestOPStackGenesis(t *testing.T) {
	for id := range superchain.OPChains {
		gen, err := LoadOPStackGenesis(id)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("chain: %d, genesis block hash: %s", id, gen.ToBlock().Hash())
	}
}
