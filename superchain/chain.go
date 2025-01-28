package superchain

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"sort"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/klauspost/compress/zstd"
)

//go:embed superchain-configs.zip
var configData []byte

var configDataReader *zip.Reader

var genesisZstdDict []byte

var Chains = make(map[uint64]*Chain)

var idsByName = make(map[string]uint64)

func ChainIDByName(name string) (uint64, error) {
	id, ok := idsByName[name]
	if !ok {
		return 0, fmt.Errorf("unknown chain %q", name)
	}
	return id, nil
}

func ChainNames() []string {
	var out []string
	for _, ch := range Chains {
		out = append(out, ch.Name+"-"+ch.Network)
	}
	sort.Strings(out)
	return out
}

func GetChain(chainID uint64) (*Chain, error) {
	chain, ok := Chains[chainID]
	if !ok {
		return nil, fmt.Errorf("unknown chain ID: %d", chainID)
	}
	return chain, nil
}

type Chain struct {
	Name    string `json:"name"`
	Network string `json:"network"`

	config  *ChainConfig
	genesis []byte

	// The config and genesis initialization is separated
	// to allow for lazy loading. Reading genesis files is
	// very expensive in Cannon so we only want to do it
	// when necessary.
	configOnce  sync.Once
	genesisOnce sync.Once
	err         error
}

func (c *Chain) Config() (*ChainConfig, error) {
	c.configOnce.Do(c.populateConfig)
	return c.config, c.err
}

func (c *Chain) GenesisData() ([]byte, error) {
	c.genesisOnce.Do(c.populateGenesis)
	return c.genesis, c.err
}

func (c *Chain) populateConfig() {
	configFile, err := configDataReader.Open(path.Join("configs", c.Network, c.Name+".toml"))
	if err != nil {
		c.err = fmt.Errorf("error opening chain config file %s/%s: %w", c.Network, c.Name, err)
		return
	}
	defer configFile.Close()

	var cfg ChainConfig
	if _, err := toml.NewDecoder(configFile).Decode(&cfg); err != nil {
		c.err = fmt.Errorf("error decoding chain config file %s/%s: %w", c.Network, c.Name, err)
		return
	}
	c.config = &cfg
}

func (c *Chain) populateGenesis() {
	genesisFile, err := configDataReader.Open(path.Join("genesis", c.Network, c.Name+".json.zst"))
	if err != nil {
		c.err = fmt.Errorf("error opening compressed genesis file %s/%s: %w", c.Network, c.Name, err)
		return
	}
	defer genesisFile.Close()
	zstdR, err := zstd.NewReader(genesisFile, zstd.WithDecoderDicts(genesisZstdDict))
	if err != nil {
		c.err = fmt.Errorf("error creating zstd reader for %s/%s: %w", c.Network, c.Name, err)
		return
	}
	defer zstdR.Close()

	out, err := io.ReadAll(zstdR)
	if err != nil {
		c.err = fmt.Errorf("error reading genesis file for %s/%s: %w", c.Network, c.Name, err)
		return
	}
	c.genesis = out
}

func init() {
	var err error
	configDataReader, err = zip.NewReader(bytes.NewReader(configData), int64(len(configData)))
	if err != nil {
		panic(fmt.Errorf("opening zip reader: %w", err))
	}
	dictR, err := configDataReader.Open("dictionary")
	if err != nil {
		panic(fmt.Errorf("error opening dictionary: %w", err))
	}
	defer dictR.Close()
	genesisZstdDict, err = io.ReadAll(dictR)
	if err != nil {
		panic(fmt.Errorf("error reading dictionary: %w", err))
	}
	chainFile, err := configDataReader.Open("chains.json")
	if err != nil {
		panic(fmt.Errorf("error opening chains file: %w", err))
	}
	defer chainFile.Close()
	if err := json.NewDecoder(chainFile).Decode(&Chains); err != nil {
		panic(fmt.Errorf("error decoding chains file: %w", err))
	}

	for chainID, chain := range Chains {
		idsByName[chain.Name+"-"+chain.Network] = chainID
	}
}
