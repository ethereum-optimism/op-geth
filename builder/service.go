package builder

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	builderSpec "github.com/attestantio/go-builder-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/flashbots/go-boost-utils/bls"
	"github.com/flashbots/go-boost-utils/ssz"
	"github.com/gorilla/mux"
)

const (
	_PathGetPayload = "/eth/v1/builder/payload/{slot:[0-9]+}/{parent_hash:0x[a-fA-F0-9]+}"
)

type Service struct {
	srv     *http.Server
	builder IBuilder
}

func (s *Service) Start() error {
	if s.srv != nil {
		log.Info("Service started")
		go s.srv.ListenAndServe()
	}

	s.builder.Start()

	return nil
}

func (s *Service) Stop() error {
	if s.srv != nil {
		s.srv.Close()
	}
	s.builder.Stop()
	return nil
}

func (s *Service) GetPayloadV1(request PayloadRequestV1) (*builderSpec.VersionedSubmitBlockRequest, error) {
	return s.builder.GetPayload(request)
}

func NewService(listenAddr string, builder IBuilder) *Service {
	var srv *http.Server

	router := mux.NewRouter()
	router.HandleFunc(_PathGetPayload, builder.handleGetPayload).Methods(http.MethodGet)

	srv = &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}

	return &Service{
		srv:     srv,
		builder: builder,
	}
}

func Register(stack *node.Node, backend *eth.Ethereum, cfg *Config) error {
	envBuilderSkBytes, err := hexutil.Decode(cfg.BuilderSecretKey)
	if err != nil {
		return errors.New("incorrect builder API secret key provided")
	}

	genesisForkVersionBytes, err := hexutil.Decode(cfg.GenesisForkVersion)
	if err != nil {
		return fmt.Errorf("invalid genesisForkVersion: %w", err)
	}

	var genesisForkVersion [4]byte
	copy(genesisForkVersion[:], genesisForkVersionBytes[:4])
	builderSigningDomain := ssz.ComputeDomain(ssz.DomainTypeAppBuilder, genesisForkVersion, phase0.Root{})

	var beaconClient IBeaconClient
	if len(cfg.BeaconEndpoints) == 0 {
		beaconClient = &NilBeaconClient{}
	} else {
		beaconClient = NewOpBeaconClient(cfg.BeaconEndpoints[0])
	}

	ethereumService := NewEthereumService(backend)

	builderSk, err := bls.SecretKeyFromBytes(envBuilderSkBytes[:])
	if err != nil {
		return errors.New("incorrect builder API secret key provided")
	}

	var builderRetryInterval time.Duration
	if cfg.RetryInterval != "" {
		d, err := time.ParseDuration(cfg.RetryInterval)
		if err != nil {
			return fmt.Errorf("error parsing builder retry interval - %v", err)
		}
		builderRetryInterval = d
	} else {
		builderRetryInterval = RetryIntervalDefault
	}

	builderArgs := BuilderArgs{
		sk:                            builderSk,
		eth:                           ethereumService,
		builderSigningDomain:          builderSigningDomain,
		builderRetryInterval:          builderRetryInterval,
		ignoreLatePayloadAttributes:   cfg.IgnoreLatePayloadAttributes,
		beaconClient:                  beaconClient,
		blockTime:                     cfg.BlockTime,
	}

	builderBackend, err := NewBuilder(builderArgs)
	if err != nil {
		return fmt.Errorf("failed to create builder backend: %w", err)
	}
	builderService := NewService(cfg.ListenAddr, builderBackend)

	stack.RegisterAPIs([]rpc.API{
		{
			Namespace:     "builder",
			Version:       "1.0",
			Service:       builderService,
			Public:        true,
			Authenticated: true,
		},
	})

	stack.RegisterLifecycle(builderService)

	return nil
}
