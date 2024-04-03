package celo

import _ "embed"

//go:generate go run ../../cmd/abigen --pkg abigen --out abigen/GoldToken.go --abi compiled/GoldToken.abi --type GoldToken
//go:generate go run ../../cmd/abigen --pkg abigen --out abigen/Registry.go --abi compiled/CeloRegistry.abi --type Registry
//go:generate go run ../../cmd/abigen --pkg abigen --out abigen/SortedOracles.go --abi compiled/SortedOracles.abi --type SortedOracles
//go:generate go run ../../cmd/abigen --pkg abigen --out abigen/FeeCurrencyWhitelist.go --abi compiled/FeeCurrencyWhitelist.abi --type FeeCurrencyWhitelist
//go:generate go run ../../cmd/abigen --pkg abigen --out abigen/FeeCurrency.go --abi compiled/FeeCurrency.abi --type FeeCurrency
//go:generate go run ../../cmd/abigen --pkg abigen --out abigen/MockSortedOracles.go --abi compiled/MockSortedOracles.abi --type MockSortedOracles

//go:embed compiled/CeloRegistry.bin-runtime
var RegistryBytecodeRaw []byte

//go:embed compiled/GoldToken.bin-runtime
var GoldTokenBytecodeRaw []byte

//go:embed compiled/Proxy.bin-runtime
var ProxyBytecodeRaw []byte

//go:embed compiled/SortedOracles.bin-runtime
var SortedOraclesBytecodeRaw []byte

//go:embed compiled/MockSortedOracles.bin-runtime
var MockSortedOraclesBytecodeRaw []byte

//go:embed compiled/FeeCurrencyWhitelist.bin-runtime
var FeeCurrencyWhitelistBytecodeRaw []byte

//go:embed compiled/FeeCurrency.bin-runtime
var FeeCurrencyBytecodeRaw []byte
