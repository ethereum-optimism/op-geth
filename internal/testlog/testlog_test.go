package testlog

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/log"
)

type mockT struct {
	out io.Writer
}

func (t *mockT) Helper() {
	// noop for the purposes of unit tests
}

func (t *mockT) Logf(format string, args ...any) {
	// we could gate this operation in a mutex, but because testlogger
	// only calls Logf with its internal mutex held, we just write output here
	var lineBuf bytes.Buffer
	if _, err := fmt.Fprintf(&lineBuf, format, args...); err != nil {
		panic(err)
	}
	// The timestamp is locale-dependent, so we want to trim that off
	// "INFO [01-01|00:00:00.000] a message ..." -> "a message..."
	sanitized := strings.Split(lineBuf.String(), "]")[1]
	if _, err := t.out.Write([]byte(sanitized)); err != nil {
		panic(err)
	}
}

func (t *mockT) FailNow() {
	panic("mock FailNow")
}

func TestLogging(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		run      func(t *mockT)
	}{
		{
			"SubLogger",
			` Visible                                  
 Hide and seek                             foobar=123
 Also visible                             
`,
			func(t *mockT) {
				l := Logger(t, log.LevelInfo)
				subLogger := l.New("foobar", 123)

				l.Info("Visible")
				subLogger.Info("Hide and seek")
				l.Info("Also visible")
			},
		},
	}

	for _, tc := range tests {
		outp := bytes.Buffer{}
		mock := mockT{&outp}
		tc.run(&mock)
		if outp.String() != tc.expected {
			fmt.Printf("output mismatch.\nwant: '%s'\ngot: '%s'\n", tc.expected, outp.String())
		}
	}
}

// TestCrit tests that Crit does not os.Exit in testing, and instead fails the test,
// while properly flushing the log output of the critical log.
func TestCrit(t *testing.T) {
	outp := bytes.Buffer{}
	mt := &mockT{&outp}
	l := Logger(mt, log.LevelInfo)
	l.Info("hello world", "a", 1)
	require.PanicsWithValue(t, "mock FailNow", func() {
		l.Crit("bye", "b", 2)
	})
	got := outp.String()
	expected := ` hello world                               a=1
 bye                                       b=2
`
	require.Equal(t, expected, got)
}
