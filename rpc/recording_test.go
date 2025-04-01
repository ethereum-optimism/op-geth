package rpc

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/internal/testlog"
	"github.com/ethereum/go-ethereum/log"
)

type testRecorder struct {
	t   require.TestingT
	log log.Logger

	inReq   int // number of incoming requests
	inResp  int // response count to incoming requests
	inNotif int // number of incoming notifications

	outReq   int // number of outgoing requests
	outResp  int // response count of outgoing requests
	outNotif int // number of sent notifications
}

var _ Recorder = (*testRecorder)(nil)

func newTestRecorder(t *testing.T, log log.Logger) *testRecorder {
	return &testRecorder{
		t:   t,
		log: log,
	}
}

func (tr *testRecorder) RecordIncoming(ctx1 context.Context, msg RecordedMsg) RecordDone {
	tr.log.Info("Handling incoming message", "id", msg.MsgID(),
		"method", msg.MsgMethod(), "notif", msg.MsgIsNotification())
	require.False(tr.t, msg.MsgIsResponse())
	if msg.MsgIsNotification() {
		tr.inNotif += 1
	} else {
		tr.inReq += 1
	}
	return func(ctx2 context.Context, input, output RecordedMsg) {
		require.Equal(tr.t, ctx1, ctx2)
		require.Equal(tr.t, msg, input)
		if output == nil {
			require.True(tr.t, input.MsgIsNotification())
			tr.log.Info("Received notification", "method", input.MsgMethod())
		} else if output != nil {
			require.True(tr.t, output.MsgIsResponse())
			require.Equal(tr.t, input.MsgID(), output.MsgID())
			tr.inResp += 1
			tr.log.Info("Received response", "id", input.MsgID(), "errObj", output.MsgError())
		}
	}
}

func (tr *testRecorder) RecordOutgoing(ctx1 context.Context, msg RecordedMsg) RecordDone {
	tr.log.Info("Handling outgoing message", "id", msg.MsgID(), "method", msg.MsgMethod())
	require.False(tr.t, msg.MsgIsResponse())
	tr.outReq += 1
	return func(ctx2 context.Context, input, output RecordedMsg) {
		tr.outResp += 1
		require.Equal(tr.t, msg, input)
		require.Equal(tr.t, ctx1, ctx2)
		if output == nil {
			require.True(tr.t, input.MsgIsNotification())
			tr.log.Info("Sent notification", "method", msg.MsgMethod())
			tr.outNotif += 1
		} else {
			require.NotNil(tr.t, output)
			require.False(tr.t, input.MsgIsResponse())
			require.True(tr.t, output.MsgIsResponse())
			require.Equal(tr.t, input.MsgID(), output.MsgID())
			tr.log.Info("Received response", "id", input.MsgID(), "errObj", output.MsgError())
		}
	}
}

func TestRecording(t *testing.T) {
	t.Parallel()

	server := newTestServer()
	defer server.Stop()
	logger := testlog.Logger(t, log.LevelError)
	srvRec := newTestRecorder(t, logger)
	server.SetRecorder(srvRec)
	clRec := newTestRecorder(t, logger)
	client := DialInProc(server, WithRecorder(clRec))
	defer client.Close()

	var resp echoResult
	if err := client.Call(&resp, "test_echo", "hello", 10, &echoArgs{"world"}); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(resp, echoResult{"hello", 10, &echoArgs{"world"}}) {
		t.Errorf("incorrect result %#v", resp)
	}
	require.Equal(t, 1, clRec.outReq)
	require.Equal(t, 1, clRec.outResp)
	require.Equal(t, 0, clRec.inReq)
	require.Equal(t, 0, clRec.inNotif)

	require.Equal(t, 0, srvRec.outReq)
	require.Equal(t, 0, srvRec.outResp)
	require.Equal(t, 1, srvRec.inReq)
	require.Equal(t, 0, srvRec.inNotif)

	nc := make(chan int)
	count := 10
	sub, err := client.Subscribe(context.Background(), "nftest", nc, "someSubscription", count, 0)
	if err != nil {
		t.Fatal("can't subscribe:", err)
	}
	for i := 0; i < count; i++ {
		if val := <-nc; val != i {
			t.Fatalf("value mismatch: got %d, want %d", val, i)
		}
	}

	sub.Unsubscribe()
	select {
	case v := <-nc:
		t.Fatal("received value after unsubscribe:", v)
	case err := <-sub.Err():
		if err != nil {
			t.Fatalf("Err returned a non-nil error after explicit unsubscribe: %q", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("subscription not closed within 1s after unsubscribe")
	}

	require.Equal(t, 10, srvRec.outNotif, "must have sent 10 notifications")
	require.Equal(t, 10, clRec.inNotif, "must have received 10 notifications")
}
