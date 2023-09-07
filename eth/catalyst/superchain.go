package catalyst

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/params"
)

var (
	requiredProtocolDeltaGauge    = metrics.NewRegisteredGauge("superchain/required/delta", nil)
	recommendedProtocolDeltaGauge = metrics.NewRegisteredGauge("superchain/recommended/delta", nil)
)

type SuperchainSignal struct {
	Recommended params.ProtocolVersion `json:"recommended"`
	Required    params.ProtocolVersion `json:"required"`
}

func (api *ConsensusAPI) SignalSuperchainV1(signal *SuperchainSignal) params.ProtocolVersion {
	if signal == nil {
		log.Info("received empty superchain version signal", "local", params.OPStackSupport)
		return params.OPStackSupport
	}
	logger := log.New("local", params.OPStackSupport, "required", signal.Required, "recommended", signal.Recommended)
	requiredCmp := params.OPStackSupport.Compare(signal.Required)
	requiredProtocolDeltaGauge.Update(int64(requiredCmp))
	switch requiredCmp {
	case params.AheadMajor:
		logger.Info("node is ahead of major required protocol change")
	case params.AheadMinor, params.AheadPatch, params.AheadPrerelease:
		logger.Debug("node is ahead of compatible required protocol change")
	case params.Matching:
		logger.Debug("node supports latest required protocol change")
	case params.OutdatedMajor:
		logger.Error("node does not support major required protocol change")
	case params.OutdatedMinor:
		logger.Warn("node does not support minor required protocol change")
	case params.OutdatedPatch:
		logger.Warn("node does not support backwards-compatible required protocol change")
	case params.OutdatedPrerelease:
		logger.Debug("new required protocol pre-release is available")
	case params.DiffBuild:
		logger.Debug("ignoring required-protocol-version signal, build is different")
	case params.DiffVersionType:
		logger.Warn("unrecognized required-protocol-version signal version-type")
	}
	recommendedCmp := params.OPStackSupport.Compare(signal.Recommended)
	recommendedProtocolDeltaGauge.Update(int64(recommendedCmp))
	switch recommendedCmp {
	case params.AheadMajor:
		logger.Info("node is ahead of major recommended protocol change")
	case params.AheadMinor, params.AheadPatch, params.AheadPrerelease:
		logger.Debug("node is ahead of compatible recommended protocol change")
	case params.Matching:
		logger.Debug("node supports latest recommended protocol change")
	case params.OutdatedMajor:
		logger.Warn("node does not support major recommended protocol change")
	case params.OutdatedMinor:
		logger.Info("node does not support minor recommended protocol change")
	case params.OutdatedPatch:
		logger.Debug("node does not support backwards-compatible recommended protocol change")
	case params.OutdatedPrerelease:
		logger.Debug("new recommended protocol pre-release is available")
	case params.DiffBuild:
		logger.Debug("ignoring recommended-protocol-version signal, build is different")
	case params.DiffVersionType:
		logger.Warn("unrecognized recommended-protocol-version signal version-type")
	}
	return params.OPStackSupport
}
