package params

import (
	"fmt"
	"testing"
)

type HumanProtocolVersion struct {
	VersionType         uint8
	Major, Minor, Patch uint32
	Prerelease          uint32
	Build               uint64
}

type ComparisonCase struct {
	A, B HumanProtocolVersion
	Cmp  ProtocolVersionComparison
}

func TestProtocolVersion_Compare(t *testing.T) {
	testCases := []ComparisonCase{
		{
			A:   HumanProtocolVersion{0, 2, 1, 1, 1, 0},
			B:   HumanProtocolVersion{0, 1, 2, 2, 2, 0},
			Cmp: AheadMajor,
		},
		{
			A:   HumanProtocolVersion{0, 1, 2, 1, 1, 0},
			B:   HumanProtocolVersion{0, 1, 1, 2, 2, 0},
			Cmp: AheadMinor,
		},
		{
			A:   HumanProtocolVersion{0, 1, 1, 2, 1, 0},
			B:   HumanProtocolVersion{0, 1, 1, 1, 2, 0},
			Cmp: AheadPatch,
		},
		{
			A:   HumanProtocolVersion{0, 1, 1, 1, 2, 0},
			B:   HumanProtocolVersion{0, 1, 1, 1, 1, 0},
			Cmp: AheadPrerelease,
		},
		{
			A:   HumanProtocolVersion{0, 1, 2, 3, 4, 0},
			B:   HumanProtocolVersion{0, 1, 2, 3, 4, 0},
			Cmp: Matching,
		},
		{
			A:   HumanProtocolVersion{0, 3, 2, 1, 5, 3},
			B:   HumanProtocolVersion{1, 1, 2, 3, 3, 6},
			Cmp: DiffVersionType,
		},
		{
			A:   HumanProtocolVersion{0, 3, 2, 1, 5, 3},
			B:   HumanProtocolVersion{0, 1, 2, 3, 3, 6},
			Cmp: DiffBuild,
		},
		{
			A:   HumanProtocolVersion{0, 0, 0, 0, 0, 0},
			B:   HumanProtocolVersion{0, 1, 3, 3, 3, 3},
			Cmp: EmptyVersion,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			a := ToProtocolVersion(tc.A.Build, tc.A.Major, tc.A.Minor, tc.A.Patch, tc.A.Prerelease)
			a[0] = tc.A.VersionType
			b := ToProtocolVersion(tc.B.Build, tc.B.Major, tc.B.Minor, tc.B.Patch, tc.B.Prerelease)
			b[0] = tc.B.VersionType
			cmp := a.Compare(b)
			if cmp != tc.Cmp {
				t.Fatalf("expected %d but got %d", tc.Cmp, cmp)
			}
			switch tc.Cmp {
			case AheadMajor, AheadMinor, AheadPatch, AheadPrerelease:
				inv := b.Compare(a)
				if inv != -tc.Cmp {
					t.Fatalf("expected inverse when reversing the comparison, %d but got %d", -tc.Cmp, inv)
				}
			case DiffVersionType, DiffBuild, EmptyVersion, Matching:
				inv := b.Compare(a)
				if inv != tc.Cmp {
					t.Fatalf("expected comparison reversed to hold the same, expected %d but got %d", tc.Cmp, inv)
				}
			}
		})
	}
}
