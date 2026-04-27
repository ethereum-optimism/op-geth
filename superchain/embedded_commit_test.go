package superchain

import (
	"os"
	"strings"
	"testing"
)

// TestEmbeddedRegistryCommit asserts that the COMMIT entry inside the embedded
// superchain-configs.zip matches the contents of superchain-registry-commit.txt
// at the repo root. They get out of sync if a developer bumps the txt file
// without re-running sync-superchain.sh.
func TestEmbeddedRegistryCommit(t *testing.T) {
	got := EmbeddedRegistryCommit()

	raw, err := os.ReadFile("../superchain-registry-commit.txt")
	if err != nil {
		t.Fatal(err)
	}
	want := strings.TrimSpace(string(raw))

	if got != want {
		t.Fatalf("embedded commit %q does not match superchain-registry-commit.txt %q;\n"+
			"the bundle is stale — run ./sync-superchain.sh from the repo root", got, want)
	}
}
