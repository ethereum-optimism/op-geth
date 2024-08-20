package builder

import (
	"net/http/httptest"

	"github.com/gorilla/mux"
)

type testBeaconClient struct {
	slot      uint64
}

func (b *testBeaconClient) Stop() {}

func (b *testBeaconClient) SubscribeToPayloadAttributesEvents(payloadAttrC chan BuilderPayloadAttributes) {
}

func (b *testBeaconClient) Start() error { return nil }

type mockBeaconNode struct {
	srv *httptest.Server
}

func newMockBeaconNode() *mockBeaconNode {
	r := mux.NewRouter()
	srv := httptest.NewServer(r)

	mbn := &mockBeaconNode{
		srv: srv,
	}

	return mbn
}
