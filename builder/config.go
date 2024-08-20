package builder

import "time"

const (
	RetryIntervalDefault = 500 * time.Millisecond
	BlockTimeDefault     = 2000 * time.Millisecond // TODO: Configure by flag.
)

type Config struct {
	Enabled                     bool          `toml:",omitempty"`
	IgnoreLatePayloadAttributes bool          `toml:",omitempty"`
	BuilderSecretKey            string        `toml:",omitempty"`
	ListenAddr                  string        `toml:",omitempty"`
	GenesisForkVersion          string        `toml:",omitempty"`
	BeaconEndpoints             []string      `toml:",omitempty"`
	RetryInterval               string        `toml:",omitempty"`
	BlockTime                   time.Duration `toml:",omitempty"`
}

// DefaultConfig is the default config for the builder.
var DefaultConfig = Config{
	Enabled:                     false,
	IgnoreLatePayloadAttributes: false,
	BuilderSecretKey:            "0x2fc12ae741f29701f8e30f5de6350766c020cb80768a0ff01e6838ffd2431e11",
	ListenAddr:                  ":28545",
	GenesisForkVersion:          "0x00000000",
	BeaconEndpoints:             []string{"http://127.0.0.1:5052"},
	RetryInterval:               RetryIntervalDefault.String(),
	BlockTime:                   BlockTimeDefault,
}
