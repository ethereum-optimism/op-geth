package catalyst

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
)

func TestSignalSuperchainV1(t *testing.T) {
	genesis, preMergeBlocks := generateMergeChain(2, false)
	n, ethservice := startEthService(t, genesis, preMergeBlocks)
	defer n.Close()
	api := NewConsensusAPI(ethservice)
	t.Run("matching", func(t *testing.T) {
		out, err := api.SignalSuperchainV1(&SuperchainSignal{
			Recommended: params.OPStackSupport,
			Required:    params.OPStackSupport,
		})
		if err != nil {
			t.Fatalf("failed to process signal: %v", err)
		}
		if out != params.OPStackSupport {
			t.Fatalf("expected %s but got %s", params.OPStackSupport, out)
		}
	})
	t.Run("null_arg", func(t *testing.T) {
		out, err := api.SignalSuperchainV1(nil)
		if err != nil {
			t.Fatalf("failed to process signal: %v", err)
		}
		if out != params.OPStackSupport {
			t.Fatalf("expected %s but got %s", params.OPStackSupport, out)
		}
	})
}

func TestSignalSuperchainV1Halt(t *testing.T) {
	// Note: depending on the params.OPStackSupport some types of bumps are not possible with active prerelease
	testCases := []struct {
		cfg  string
		bump string
		halt bool
	}{
		{"none", "major", false},
		{"major", "major", true},
		{"minor", "major", true},
		{"patch", "major", true},
		{"major", "minor", false},
		{"minor", "minor", true},
		{"patch", "minor", true},
		{"major", "patch", false},
		{"minor", "patch", false},
		{"patch", "patch", true},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.cfg+"_"+tc.bump, func(t *testing.T) {
			genesis, preMergeBlocks := generateMergeChain(2, false)
			ethcfg := &ethconfig.Config{Genesis: genesis, SyncMode: downloader.FullSync, TrieTimeout: time.Minute, TrieDirtyCache: 256, TrieCleanCache: 256}
			ethcfg.RollupHaltOnIncompatibleProtocolVersion = tc.cfg // opt-in to halting (or not)
			n, ethservice := startEthServiceWithConfigFn(t, preMergeBlocks, ethcfg)
			defer n.Close() // close at the end, regardless of any prior (failed) closing
			api := NewConsensusAPI(ethservice)
			_, build, major, minor, patch, preRelease := params.OPStackSupport.Parse()
			if preRelease != 0 { // transform back from prerelease, so we can do a clean version bump
				if patch != 0 {
					patch -= 1
				} else if minor != 0 {
					// can't patch-bump e.g. v3.1.0-1, the prerelease forces a minor bump:
					// v3.0.999 is lower than the prerelease, v3.1.1-1 is a prerelease of v3.1.1.
					if tc.bump == "patch" {
						t.Skip()
					}
					minor -= 1
				} else if major != 0 {
					if tc.bump == "minor" || tc.bump == "patch" { // can't minor-bump or patch-bump
						t.Skip()
					}
					major -= 1
				}
				preRelease = 0
			}
			majorSignal, minorSignal, patchSignal := major, minor, patch
			switch tc.bump {
			case "major":
				majorSignal += 1
			case "minor":
				minorSignal += 1
			case "patch":
				patchSignal += 1
			}
			out, err := api.SignalSuperchainV1(&SuperchainSignal{
				Recommended: params.OPStackSupport, // required version change should be enough
				Required:    params.ProtocolVersionV0{Build: build, Major: majorSignal, Minor: minorSignal, Patch: patchSignal, PreRelease: preRelease}.Encode(),
			})
			if err != nil {
				t.Fatalf("failed to process signal: %v", err)
			}
			if out != params.OPStackSupport {
				t.Fatalf("expected %s but got %s", params.OPStackSupport, out)
			}
			closeErr := n.Close()
			if !tc.halt {
				// assert no halt by closing, and not getting any error
				if closeErr != nil {
					t.Fatalf("expected not to have closed already, but just closed without error")
				}
			} else {
				// assert halt by closing again, and seeing if we are not stopped already
				if closeErr != node.ErrNodeStopped {
					t.Fatalf("expected to have already closed and get a ErrNodeStopped error, but got %v", closeErr)
				}
			}
		})
	}
}
