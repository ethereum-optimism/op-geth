package rpc

import (
	"context"
	"encoding/json"
)

// RecordDone is called after an incoming input (request or notification) has successfully been processed,
// and captures the output (nil or response).
type RecordDone func(ctx context.Context, input, output RecordedMsg)

// Recorder captures RPC traffic
type Recorder interface {
	// RecordIncoming records an incoming message (request or notification), before it has been processed.
	// It may optionally return a function to capture the result of the processing (response or nil).
	RecordIncoming(ctx context.Context, msg RecordedMsg) RecordDone

	// RecordOutgoing records an outgoing message (request or notification), before it has been sent.
	// It may optionally return a function to capture the result of the request (response or nil).
	RecordOutgoing(ctx context.Context, msg RecordedMsg) RecordDone
}

// RecordedMsg wraps around the internal jsonrpcMessage type,
// to provide a public read-only interface for recording of RPC activity.
// Every method name is prefixed with "Msg", to avoid conflict with internal methods and future geth changes.
type RecordedMsg interface {
	MsgIsNotification() bool
	MsgIsResponse() bool
	MsgID() json.RawMessage
	MsgMethod() string
	MsgParams() json.RawMessage
	MsgError() *JsonError
	MsgResult() json.RawMessage
}

var _ RecordedMsg = (*jsonrpcMessage)(nil)

func (msg *jsonrpcMessage) MsgIsNotification() bool {
	return msg.isNotification()
}

func (msg *jsonrpcMessage) MsgIsResponse() bool {
	return msg.isResponse()
}

func (msg *jsonrpcMessage) MsgID() json.RawMessage {
	return msg.ID
}
func (msg *jsonrpcMessage) MsgMethod() string {
	return msg.Method
}
func (msg *jsonrpcMessage) MsgParams() json.RawMessage {
	return msg.Params
}
func (msg *jsonrpcMessage) MsgError() *JsonError {
	return msg.Error
}
func (msg *jsonrpcMessage) MsgResult() json.RawMessage {
	return msg.Result
}

var _ RecordedMsg = (*jsonrpcSubscriptionNotification)(nil)

func (notif *jsonrpcSubscriptionNotification) MsgIsNotification() bool {
	return true
}

func (notif *jsonrpcSubscriptionNotification) MsgIsResponse() bool {
	return false
}

func (notif *jsonrpcSubscriptionNotification) MsgID() json.RawMessage {
	return json.RawMessage{} // notifications do not have an ID
}

func (notif *jsonrpcSubscriptionNotification) MsgMethod() string {
	return notif.Method
}

func (notif *jsonrpcSubscriptionNotification) MsgParams() json.RawMessage {
	data, _ := json.Marshal(notif.Params)
	return data
}

func (notif *jsonrpcSubscriptionNotification) MsgError() *JsonError {
	return nil
}

func (notif *jsonrpcSubscriptionNotification) MsgResult() json.RawMessage {
	return nil
}

// WithRecorder attaches a recorder to an RPC client, useful when it is serving bidirectional RPC
func WithRecorder(r Recorder) ClientOption {
	return optionFunc(func(cfg *clientConfig) {
		cfg.recorder = r
	})
}
