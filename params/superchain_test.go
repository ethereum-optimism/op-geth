package params

import (
	"fmt"
	"testing"
)

type HumanProtocolVersion struct {
	VersionType         uint8
	Major, Minor, Patch uint32
	Prerelease          uint32
	Build               [8]byte
}

type ComparisonCase struct {
	A, B HumanProtocolVersion
	Cmp  ProtocolVersionComparison
}

func TestProtocolVersion_Compare(t *testing.T) {
	testCases := []ComparisonCase{
		{
			A:   HumanProtocolVersion{0, 2, 1, 1, 1, [8]byte{}},
			B:   HumanProtocolVersion{0, 1, 2, 2, 2, [8]byte{}},
			Cmp: AheadMajor,
		},
		{
			A:   HumanProtocolVersion{0, 1, 2, 1, 1, [8]byte{}},
			B:   HumanProtocolVersion{0, 1, 1, 2, 2, [8]byte{}},
			Cmp: AheadMinor,
		},
		{
			A:   HumanProtocolVersion{0, 1, 1, 2, 1, [8]byte{}},
			B:   HumanProtocolVersion{0, 1, 1, 1, 2, [8]byte{}},
			Cmp: AheadPatch,
		},
		{
			A:   HumanProtocolVersion{0, 1, 1, 1, 2, [8]byte{}},
			B:   HumanProtocolVersion{0, 1, 1, 1, 1, [8]byte{}},
			Cmp: AheadPrerelease,
		},
		{
			A:   HumanProtocolVersion{0, 1, 2, 3, 4, [8]byte{}},
			B:   HumanProtocolVersion{0, 1, 2, 3, 4, [8]byte{}},
			Cmp: Matching,
		},
		{
			A:   HumanProtocolVersion{0, 3, 2, 1, 5, [8]byte{3}},
			B:   HumanProtocolVersion{1, 1, 2, 3, 3, [8]byte{6}},
			Cmp: DiffVersionType,
		},
		{
			A:   HumanProtocolVersion{0, 3, 2, 1, 5, [8]byte{3}},
			B:   HumanProtocolVersion{0, 1, 2, 3, 3, [8]byte{6}},
			Cmp: DiffBuild,
		},
		{
			A:   HumanProtocolVersion{0, 0, 0, 0, 0, [8]byte{}},
			B:   HumanProtocolVersion{0, 1, 3, 3, 3, [8]byte{3}},
			Cmp: EmptyVersion,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			a := ProtocolVersionV0{tc.A.Build, tc.A.Major, tc.A.Minor, tc.A.Patch, tc.A.Prerelease}.Encode()
			a[0] = tc.A.VersionType
			b := ProtocolVersionV0{tc.B.Build, tc.B.Major, tc.B.Minor, tc.B.Patch, tc.B.Prerelease}.Encode()
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
func TestProtocolVersion_String(t *testing.T) {
	testCases := []struct {
		version  ProtocolVersion
		expected string
	}{
		{ProtocolVersionV0{[8]byte{}, 0, 0, 0, 0}.Encode(), "v0.0.0"},
		{ProtocolVersionV0{[8]byte{}, 0, 0, 0, 1}.Encode(), "v0.0.0-1"},
		{ProtocolVersionV0{[8]byte{}, 0, 0, 1, 0}.Encode(), "v0.0.1"},
		{ProtocolVersionV0{[8]byte{}, 4, 3, 2, 1}.Encode(), "v4.3.2-1"},
		{ProtocolVersionV0{[8]byte{}, 0, 100, 2, 0}.Encode(), "v0.100.2"},
		{ProtocolVersionV0{[8]byte{'O', 'P', '-', 'm', 'o', 'd'}, 42, 0, 2, 1}.Encode(), "v42.0.2-1+OP-mod"},
		{ProtocolVersionV0{[8]byte{'b', 'e', 't', 'a', '.', '1', '2', '3'}, 1, 0, 0, 0}.Encode(), "v1.0.0+beta.123"},
		{ProtocolVersionV0{[8]byte{'a', 'b', 1}, 42, 0, 2, 0}.Encode(), "v42.0.2+0x6162010000000000"}, // do not render invalid alpha numeric
		{ProtocolVersionV0{[8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 42, 0, 2, 0}.Encode(), "v42.0.2+0x0102030405060708"},
	}
	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			got := tc.version.String()
			if got != tc.expected {
				t.Fatalf("got %q but expected %q", got, tc.expected)
			}
		})
	}
}
