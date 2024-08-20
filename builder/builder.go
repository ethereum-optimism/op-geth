package builder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	_ "os"
	"strconv"
	"sync"
	"time"

	builderApi "github.com/attestantio/go-builder-client/api"
	builderApiBellatrix "github.com/attestantio/go-builder-client/api/bellatrix"
	builderApiCapella "github.com/attestantio/go-builder-client/api/capella"
	builderApiDeneb "github.com/attestantio/go-builder-client/api/deneb"
	builderApiV1 "github.com/attestantio/go-builder-client/api/v1"
	builderSpec "github.com/attestantio/go-builder-client/spec"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/flashbots/go-boost-utils/bls"
	"github.com/flashbots/go-boost-utils/ssz"
	boostTypes "github.com/flashbots/go-boost-utils/types"
	"github.com/flashbots/go-boost-utils/utils"
	"github.com/gorilla/mux"
	"github.com/holiman/uint256"
)

var (
	ErrIncorrectSlot         = errors.New("incorrect slot")
	ErrNoPayloads            = errors.New("no payloads")
	ErrSlotFromPayload       = errors.New("could not get slot from payload")
	ErrSlotMismatch          = errors.New("slot mismatch")
	ErrParentHashFromPayload = errors.New("could not get parent hash from payload")
	ErrParentHashMismatch    = errors.New("parent hash mismatch")
)

type IBuilder interface {
	GetPayload(request PayloadRequestV1) (*builderSpec.VersionedSubmitBlockRequest, error)
	Start() error
	Stop() error

	handleGetPayload(w http.ResponseWriter, req *http.Request)
}

type Builder struct {
	eth                         IEthereumService
	ignoreLatePayloadAttributes bool
	beaconClient                IBeaconClient
	builderSecretKey            *bls.SecretKey
	builderPublicKey            phase0.BLSPubKey
	builderSigningDomain        phase0.Domain

	builderRetryInterval          time.Duration
	builderBlockTime              time.Duration
	submissionOffsetFromEndOfSlot time.Duration

	slotMu        sync.Mutex
	slotAttrs     BuilderPayloadAttributes
	slotCtx       context.Context
	slotCtxCancel context.CancelFunc

	bestBlockMu sync.Mutex
	bestBlock   *builderSpec.VersionedSubmitBlockRequest

	stop chan struct{}
}

// BuilderArgs is a struct that contains all the arguments needed to create a new Builder
type BuilderArgs struct {
	sk                            *bls.SecretKey
	builderSigningDomain          phase0.Domain
	builderRetryInterval          time.Duration
	blockTime                     time.Duration
	eth                           IEthereumService
	ignoreLatePayloadAttributes   bool
	beaconClient                  IBeaconClient
	submissionOffsetFromEndOfSlot time.Duration
}

// SubmitBlockOpts is a struct that contains all the arguments needed to submit a block to the relay
type SubmitBlockOpts struct {
	// ExecutablePayloadEnvelope is the payload envelope that was executed
	ExecutionPayloadEnvelope *engine.ExecutionPayloadEnvelope
	// SealedAt is the time at which the block was sealed
	SealedAt time.Time
	// ProposerPubkey is the proposer's pubkey
	ProposerPubkey phase0.BLSPubKey
	// PayloadAttributes are the payload attributes used for block building
	PayloadAttributes *BuilderPayloadAttributes
}

func NewBuilder(args BuilderArgs) (*Builder, error) {
	blsPk, err := bls.PublicKeyFromSecretKey(args.sk)
	if err != nil {
		return nil, err
	}
	pk, err := utils.BlsPublicKeyToPublicKey(blsPk)
	if err != nil {
		return nil, err
	}

	slotCtx, slotCtxCancel := context.WithCancel(context.Background())
	return &Builder{
		eth:                           args.eth,
		ignoreLatePayloadAttributes:   args.ignoreLatePayloadAttributes,
		beaconClient:                  args.beaconClient,
		builderSecretKey:              args.sk,
		builderPublicKey:              pk,
		builderSigningDomain:          args.builderSigningDomain,
		builderRetryInterval:          args.builderRetryInterval,
		builderBlockTime:              args.blockTime,
		submissionOffsetFromEndOfSlot: args.submissionOffsetFromEndOfSlot,

		slotCtx:       slotCtx,
		slotCtxCancel: slotCtxCancel,

		stop: make(chan struct{}, 1),
	}, nil
}

func (b *Builder) Start() error {
	log.Info("Starting builder")
	// Start regular payload attributes updates
	go func() {
		c := make(chan BuilderPayloadAttributes)
		go b.beaconClient.SubscribeToPayloadAttributesEvents(c)

		currentSlot := uint64(0)

		for {
			select {
			case <-b.stop:
				return
			case payloadAttributes := <-c:
				log.Info("Received payload attributes", "slot", payloadAttributes.Slot, "hash", payloadAttributes.HeadHash.String())
				// Right now we are building only on a single head. This might change in the future!
				if payloadAttributes.Slot < currentSlot {
					continue
				} else if payloadAttributes.Slot == currentSlot {
					// Subsequent sse events should only be canonical!
					if !b.ignoreLatePayloadAttributes {
						err := b.handlePayloadAttributes(&payloadAttributes)
						if err != nil {
							log.Error("error with builder processing on payload attribute",
								"latestSlot", currentSlot,
								"processedSlot", payloadAttributes.Slot,
								"headHash", payloadAttributes.HeadHash.String(),
								"error", err)
						}
					}
				} else if payloadAttributes.Slot > currentSlot {
					currentSlot = payloadAttributes.Slot
					err := b.handlePayloadAttributes(&payloadAttributes)
					if err != nil {
						log.Error("error with builder processing on payload attribute",
							"latestSlot", currentSlot,
							"processedSlot", payloadAttributes.Slot,
							"headHash", payloadAttributes.HeadHash.String(),
							"error", err)
					}
				}
			}
		}
	}()

	return b.beaconClient.Start()
}

func (b *Builder) Stop() error {
	close(b.stop)
	return nil
}

func (b *Builder) GetPayload(request PayloadRequestV1) (*builderSpec.VersionedSubmitBlockRequest, error) {
	log.Info("received get payload request", "slot", request.Slot, "parent", request.ParentHash)
	b.bestBlockMu.Lock()
	bestBlock := b.bestBlock
	b.bestBlockMu.Unlock()

	if bestBlock == nil {
		log.Error("no builder submissions")
		return nil, ErrNoPayloads
	}

	submittedSlot, err := bestBlock.Slot()
	if err != nil {
		log.Error("could not get slot from best submission", "err", err)
		return nil, ErrSlotFromPayload
	}

	if submittedSlot != uint64(request.Slot) {
		log.Error("slot not equal", "requested", request.Slot, "block", submittedSlot)
		return nil, ErrSlotMismatch
	}

	submittedParentHash, err := bestBlock.ParentHash()
	if err != nil {
		log.Error("could not get parent hash from best submission", "err", err)
		return nil, ErrParentHashFromPayload
	}

	if submittedParentHash.String() != request.ParentHash.String() {
		log.Error("parent hash not equal", "requested", request.ParentHash, "block", submittedParentHash.String())
		return nil, ErrParentHashMismatch
	}

	blockHash, err := bestBlock.BlockHash()
	if err != nil {
		log.Warn("could not get block hash from best submission", "err", err)
	} else {
		log.Info("payload delivered", "hash", blockHash.String())
	}

	return bestBlock, nil
}

func (b *Builder) handleGetPayload(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	slot, err := strconv.Atoi(vars["slot"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "incorrect slot")
		return
	}
	parentHashHex := vars["parent_hash"]

	log.Info("received handle get payload request", "slot", slot, "parent", parentHashHex)

	bestSubmission, err := b.GetPayload(PayloadRequestV1{
		Slot:       uint64(slot),
		ParentHash: common.HexToHash(parentHashHex),
	})
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(bestSubmission); err != nil {
		log.Error("could not encode response", "err", err)
		respondError(w, http.StatusInternalServerError, "could not encode response")
		return
	}
}

func (b *Builder) saveBlockSubmission(opts SubmitBlockOpts) error {
	executionPayload := opts.ExecutionPayloadEnvelope.ExecutionPayload
	log.Info(
		"saveBlockSubmission",
		"slot", opts.PayloadAttributes.Slot,
		"parent", opts.PayloadAttributes.HeadHash.String(),
		"hash", executionPayload.BlockHash.String(),
	)

	var dataVersion spec.DataVersion
	if b.eth.Config().IsEcotone(executionPayload.Timestamp) {
		dataVersion = spec.DataVersionDeneb
	} else if b.eth.Config().IsCanyon(executionPayload.Timestamp) {
		dataVersion = spec.DataVersionCapella
	} else {
		dataVersion = spec.DataVersionBellatrix
	}

	value, overflow := uint256.FromBig(opts.ExecutionPayloadEnvelope.BlockValue)
	if overflow {
		return fmt.Errorf("could not set block value due to value overflow")
	}

	blockBidMsg := builderApiV1.BidTrace{
		Slot:                 opts.PayloadAttributes.Slot,
		ParentHash:           phase0.Hash32(executionPayload.ParentHash),
		BlockHash:            phase0.Hash32(executionPayload.BlockHash),
		BuilderPubkey:        b.builderPublicKey,
		ProposerPubkey:       opts.ProposerPubkey,
		ProposerFeeRecipient: bellatrix.ExecutionAddress(opts.PayloadAttributes.SuggestedFeeRecipient),
		GasLimit:             executionPayload.GasLimit,
		GasUsed:              executionPayload.GasUsed,
		Value:                value,
	}

	versionedBlockRequest, err := b.getBlockRequest(opts.ExecutionPayloadEnvelope, dataVersion, &blockBidMsg)
	if err != nil {
		log.Error("could not get block request", "err", err)
		return err
	}

	b.bestBlockMu.Lock()
	b.bestBlock = versionedBlockRequest
	b.bestBlockMu.Unlock()

	log.Info("saved block", "version", dataVersion.String(), "slot", opts.PayloadAttributes.Slot, "value", opts.ExecutionPayloadEnvelope.BlockValue.String(),
		"parent", executionPayload.ParentHash.String(), "hash", executionPayload.BlockHash)

	return nil
}

func (b *Builder) getBlockRequest(executableData *engine.ExecutionPayloadEnvelope, dataVersion spec.DataVersion, blockBidMsg *builderApiV1.BidTrace) (*builderSpec.VersionedSubmitBlockRequest, error) {
	payload, err := executableDataToExecutionPayload(executableData, dataVersion)
	if err != nil {
		log.Error("could not format execution payload", "err", err)
		return nil, err
	}

	signature, err := ssz.SignMessage(blockBidMsg, b.builderSigningDomain, b.builderSecretKey)
	if err != nil {
		log.Error("could not sign builder bid", "err", err)
		return nil, err
	}

	var versionedBlockRequest builderSpec.VersionedSubmitBlockRequest
	switch dataVersion {
	case spec.DataVersionBellatrix:
		blockSubmitReq := builderApiBellatrix.SubmitBlockRequest{
			Signature:        signature,
			Message:          blockBidMsg,
			ExecutionPayload: payload.Bellatrix,
		}
		versionedBlockRequest = builderSpec.VersionedSubmitBlockRequest{
			Version:   spec.DataVersionBellatrix,
			Bellatrix: &blockSubmitReq,
		}
	case spec.DataVersionCapella:
		blockSubmitReq := builderApiCapella.SubmitBlockRequest{
			Signature:        signature,
			Message:          blockBidMsg,
			ExecutionPayload: payload.Capella,
		}
		versionedBlockRequest = builderSpec.VersionedSubmitBlockRequest{
			Version: spec.DataVersionCapella,
			Capella: &blockSubmitReq,
		}
	case spec.DataVersionDeneb:
		blockSubmitReq := builderApiDeneb.SubmitBlockRequest{
			Signature:        signature,
			Message:          blockBidMsg,
			ExecutionPayload: payload.Deneb.ExecutionPayload,
			BlobsBundle:      payload.Deneb.BlobsBundle,
		}
		versionedBlockRequest = builderSpec.VersionedSubmitBlockRequest{
			Version: spec.DataVersionDeneb,
			Deneb:   &blockSubmitReq,
		}
	}
	return &versionedBlockRequest, err
}

func (b *Builder) handlePayloadAttributes(attrs *BuilderPayloadAttributes) error {
	log.Info("Payload attribute received", "slot", attrs.Slot, "hash", attrs.HeadHash, "txs", attrs.Transactions)
	if attrs == nil {
		return nil
	}

	parentBlock := b.eth.GetBlockByHash(attrs.HeadHash)
	if parentBlock == nil {
		return fmt.Errorf("parent block hash not found in block tree given head block hash %s", attrs.HeadHash)
	}

	proposerPubkey := phase0.BLSPubKey{}

	if !b.eth.Synced() {
		return errors.New("backend not Synced")
	}

	b.slotMu.Lock()
	defer b.slotMu.Unlock()

	if attrs.Equal(&b.slotAttrs) {
		log.Debug("ignoring known payload attribute", "slot", attrs.Slot, "hash", attrs.HeadHash)
		return nil
	}

	if b.slotCtxCancel != nil {
		b.slotCtxCancel()
	}

	slotCtx, slotCtxCancel := context.WithTimeout(context.Background(), b.builderBlockTime)
	b.slotAttrs = *attrs
	b.slotCtx = slotCtx
	b.slotCtxCancel = slotCtxCancel

	go b.runBuildingJob(b.slotCtx, proposerPubkey, attrs)
	return nil
}

func (b *Builder) runBuildingJob(slotCtx context.Context, proposerPubkey phase0.BLSPubKey, attrs *BuilderPayloadAttributes) {
	ctx, cancel := context.WithTimeout(slotCtx, b.builderBlockTime)
	defer cancel()

	// Submission queue for the given payload attributes
	// multiple jobs can run for different attributes fot the given slot
	// 1. When new block is ready we check if its profit is higher than profit of last best block
	//    if it is we set queueBest* to values of the new block and notify queueSignal channel.
	var (
		queueMu                sync.Mutex
		queueLastSubmittedHash common.Hash
		queueBestBlockValue    *big.Int = big.NewInt(0)
	)

	log.Info("runBuildingJob", "slot", attrs.Slot, "parent", attrs.HeadHash, "payloadTimestamp", uint64(attrs.Timestamp), "txs", attrs.Transactions)

	// retry build block every builderBlockRetryInterval
	runRetryLoop(ctx, b.builderRetryInterval, func() {
		log.Info("retrying BuildBlock",
			"slot", attrs.Slot,
			"parent", attrs.HeadHash,
			"retryInterval", b.builderRetryInterval.String())
		payload, err := b.eth.BuildBlock(attrs)
		if err != nil {
			log.Warn("Failed to build block", "err", err)
			return
		}

		sealedAt := time.Now()
		queueMu.Lock()
		defer queueMu.Unlock()
		if payload.ExecutionPayload.BlockHash != queueLastSubmittedHash && payload.BlockValue.Cmp(queueBestBlockValue) > 0 {
			queueLastSubmittedHash = payload.ExecutionPayload.BlockHash
			queueBestBlockValue = payload.BlockValue

			submitBlockOpts := SubmitBlockOpts{
				ExecutionPayloadEnvelope: payload,
				SealedAt:                 sealedAt,
				ProposerPubkey:           proposerPubkey,
				PayloadAttributes:        attrs,
			}
			err := b.saveBlockSubmission(submitBlockOpts)
			if err != nil {
				log.Error("could not save block submission", "err", err)
			}
		}
	})
}

func executableDataToExecutionPayload(data *engine.ExecutionPayloadEnvelope, version spec.DataVersion) (*builderApi.VersionedSubmitBlindedBlockResponse, error) {
	// if version in phase0, altair, unsupported version
	if version == spec.DataVersionUnknown || version == spec.DataVersionPhase0 || version == spec.DataVersionAltair {
		return nil, fmt.Errorf("unsupported data version %d", version)
	}

	payload := data.ExecutionPayload
	blobsBundle := data.BlobsBundle

	transactionData := make([]bellatrix.Transaction, len(payload.Transactions))
	for i, tx := range payload.Transactions {
		transactionData[i] = bellatrix.Transaction(tx)
	}

	baseFeePerGas := new(boostTypes.U256Str)
	err := baseFeePerGas.FromBig(payload.BaseFeePerGas)
	if err != nil {
		return nil, err
	}

	if version == spec.DataVersionBellatrix {
		return getBellatrixPayload(payload, *baseFeePerGas, transactionData), nil
	}

	withdrawalData := make([]*capella.Withdrawal, len(payload.Withdrawals))
	for i, wd := range payload.Withdrawals {
		withdrawalData[i] = &capella.Withdrawal{
			Index:          capella.WithdrawalIndex(wd.Index),
			ValidatorIndex: phase0.ValidatorIndex(wd.Validator),
			Address:        bellatrix.ExecutionAddress(wd.Address),
			Amount:         phase0.Gwei(wd.Amount),
		}
	}
	if version == spec.DataVersionCapella {
		return getCapellaPayload(payload, *baseFeePerGas, transactionData, withdrawalData), nil
	}

	uint256BaseFeePerGas, overflow := uint256.FromBig(payload.BaseFeePerGas)
	if overflow {
		return nil, fmt.Errorf("base fee per gas overflow")
	}

	if len(blobsBundle.Blobs) != len(blobsBundle.Commitments) || len(blobsBundle.Blobs) != len(blobsBundle.Proofs) {
		return nil, fmt.Errorf("blobs bundle length mismatch")
	}

	if version == spec.DataVersionDeneb {
		return getDenebPayload(payload, uint256BaseFeePerGas, transactionData, withdrawalData, blobsBundle), nil
	}

	return nil, fmt.Errorf("unsupported data version %d", version)
}

func getBellatrixPayload(
	payload *engine.ExecutableData,
	baseFeePerGas [32]byte,
	transactions []bellatrix.Transaction,
) *builderApi.VersionedSubmitBlindedBlockResponse {
	return &builderApi.VersionedSubmitBlindedBlockResponse{
		Version: spec.DataVersionBellatrix,
		Bellatrix: &bellatrix.ExecutionPayload{
			ParentHash:    [32]byte(payload.ParentHash),
			FeeRecipient:  [20]byte(payload.FeeRecipient),
			StateRoot:     [32]byte(payload.StateRoot),
			ReceiptsRoot:  [32]byte(payload.ReceiptsRoot),
			LogsBloom:     types.BytesToBloom(payload.LogsBloom),
			PrevRandao:    [32]byte(payload.Random),
			BlockNumber:   payload.Number,
			GasLimit:      payload.GasLimit,
			GasUsed:       payload.GasUsed,
			Timestamp:     payload.Timestamp,
			ExtraData:     payload.ExtraData,
			BaseFeePerGas: baseFeePerGas,
			BlockHash:     [32]byte(payload.BlockHash),
			Transactions:  transactions,
		},
	}
}

func getCapellaPayload(
	payload *engine.ExecutableData,
	baseFeePerGas [32]byte,
	transactions []bellatrix.Transaction,
	withdrawals []*capella.Withdrawal,
) *builderApi.VersionedSubmitBlindedBlockResponse {
	return &builderApi.VersionedSubmitBlindedBlockResponse{
		Version: spec.DataVersionCapella,
		Capella: &capella.ExecutionPayload{
			ParentHash:    [32]byte(payload.ParentHash),
			FeeRecipient:  [20]byte(payload.FeeRecipient),
			StateRoot:     [32]byte(payload.StateRoot),
			ReceiptsRoot:  [32]byte(payload.ReceiptsRoot),
			LogsBloom:     types.BytesToBloom(payload.LogsBloom),
			PrevRandao:    [32]byte(payload.Random),
			BlockNumber:   payload.Number,
			GasLimit:      payload.GasLimit,
			GasUsed:       payload.GasUsed,
			Timestamp:     payload.Timestamp,
			ExtraData:     payload.ExtraData,
			BaseFeePerGas: baseFeePerGas,
			BlockHash:     [32]byte(payload.BlockHash),
			Transactions:  transactions,
			Withdrawals:   withdrawals,
		},
	}
}

func getBlobsBundle(blobsBundle *engine.BlobsBundleV1) *builderApiDeneb.BlobsBundle {
	commitments := make([]deneb.KZGCommitment, len(blobsBundle.Commitments))
	proofs := make([]deneb.KZGProof, len(blobsBundle.Proofs))
	blobs := make([]deneb.Blob, len(blobsBundle.Blobs))

	// we assume the lengths for blobs bundle is validated beforehand to be the same
	for i := range blobsBundle.Blobs {
		var commitment deneb.KZGCommitment
		copy(commitment[:], blobsBundle.Commitments[i][:])
		commitments[i] = commitment

		var proof deneb.KZGProof
		copy(proof[:], blobsBundle.Proofs[i][:])
		proofs[i] = proof

		var blob deneb.Blob
		copy(blob[:], blobsBundle.Blobs[i][:])
		blobs[i] = blob
	}
	return &builderApiDeneb.BlobsBundle{
		Commitments: commitments,
		Proofs:      proofs,
		Blobs:       blobs,
	}
}

func getDenebPayload(
	payload *engine.ExecutableData,
	baseFeePerGas *uint256.Int,
	transactions []bellatrix.Transaction,
	withdrawals []*capella.Withdrawal,
	blobsBundle *engine.BlobsBundleV1,
) *builderApi.VersionedSubmitBlindedBlockResponse {
	return &builderApi.VersionedSubmitBlindedBlockResponse{
		Version: spec.DataVersionDeneb,
		Deneb: &builderApiDeneb.ExecutionPayloadAndBlobsBundle{
			ExecutionPayload: &deneb.ExecutionPayload{
				ParentHash:    [32]byte(payload.ParentHash),
				FeeRecipient:  [20]byte(payload.FeeRecipient),
				StateRoot:     [32]byte(payload.StateRoot),
				ReceiptsRoot:  [32]byte(payload.ReceiptsRoot),
				LogsBloom:     types.BytesToBloom(payload.LogsBloom),
				PrevRandao:    [32]byte(payload.Random),
				BlockNumber:   payload.Number,
				GasLimit:      payload.GasLimit,
				GasUsed:       payload.GasUsed,
				Timestamp:     payload.Timestamp,
				ExtraData:     payload.ExtraData,
				BaseFeePerGas: baseFeePerGas,
				BlockHash:     [32]byte(payload.BlockHash),
				Transactions:  transactions,
				Withdrawals:   withdrawals,
				BlobGasUsed:   *payload.BlobGasUsed,
				ExcessBlobGas: *payload.ExcessBlobGas,
			},
			BlobsBundle: getBlobsBundle(blobsBundle),
		},
	}
}
