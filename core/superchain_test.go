package core

import "testing"

func TestOPStackGenesis(t *testing.T) {
	gen, err := LoadOPStackGenesis(7777777)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("genesis block hash: %s", gen.ToBlock().Hash())
}
