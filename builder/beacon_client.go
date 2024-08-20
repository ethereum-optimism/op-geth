package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/r3labs/sse"
)

type IBeaconClient interface {
	SubscribeToPayloadAttributesEvents(payloadAttrC chan BuilderPayloadAttributes)
	Start() error
	Stop()
}

type NilBeaconClient struct{}

func (b *NilBeaconClient) SubscribeToPayloadAttributesEvents(payloadAttrC chan BuilderPayloadAttributes) {
}

func (b *NilBeaconClient) Start() error { return nil }

func (b *NilBeaconClient) Stop() {}

type OpBeaconClient struct {
	ctx      context.Context
	cancelFn context.CancelFunc

	endpoint string
}

func NewOpBeaconClient(endpoint string) *OpBeaconClient {
	log.Info("creating opbeacon client", "endpoint", endpoint)
	ctx, cancelFn := context.WithCancel(context.Background())
	return &OpBeaconClient{
		ctx:      ctx,
		cancelFn: cancelFn,

		endpoint: endpoint,
	}
}

func (opbc *OpBeaconClient) SubscribeToPayloadAttributesEvents(payloadAttrC chan BuilderPayloadAttributes) {
	eventsURL := fmt.Sprintf("%s/events", opbc.endpoint)
	log.Info("subscribing to payload_attributes events opbs", "eventsURL", eventsURL)

	for {
		client := sse.NewClient(eventsURL)
		err := client.SubscribeWithContext(opbc.ctx, "payload_attributes", func(msg *sse.Event) {
			data := new(BuilderPayloadAttributes)
			err := json.Unmarshal(msg.Data, data)
			if err != nil {
				log.Error("could not unmarshal payload_attributes event", "err", err)
			} else {
				payloadAttrC <- *data
			}
		})
		if err != nil {
			log.Error("failed to subscribe to payload_attributes events", "err", err)
			time.Sleep(1 * time.Second)
		}
		log.Warn("opnode Subscribe ended, reconnecting")
	}
}

func (opbc *OpBeaconClient) Start() error {
	return nil
}

func (opbc *OpBeaconClient) Stop() {
	opbc.cancelFn()
}
