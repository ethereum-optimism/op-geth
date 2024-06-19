package celo

import _ "embed"

//go:generate go run ../../cmd/abigen --pkg abigen --out abigen/CeloToken.go --abi compiled/GoldToken.abi --type CeloToken
//go:generate go run ../../cmd/abigen --pkg abigen --out abigen/Registry.go --abi compiled/CeloRegistry.abi --type Registry
//go:generate go run ../../cmd/abigen --pkg abigen --out abigen/FeeCurrency.go --abi compiled/FeeCurrency.abi --type FeeCurrency
//go:generate go run ../../cmd/abigen --pkg abigen --out abigen/FeeCurrencyDirectory.go --abi compiled/FeeCurrencyDirectory.abi --type FeeCurrencyDirectory
//go:generate go run ../../cmd/abigen --pkg abigen --out abigen/MockOracle.go --abi compiled/MockOracle.abi --type MockOracle

//go:embed compiled/CeloRegistry.bin-runtime
var RegistryBytecodeRaw []byte

//go:embed compiled/GoldToken.bin-runtime
var CeloTokenBytecodeRaw []byte

//go:embed compiled/Proxy.bin-runtime
var ProxyBytecodeRaw []byte

//go:embed compiled/FeeCurrency.bin-runtime
var FeeCurrencyBytecodeRaw []byte

//go:embed compiled/FeeCurrencyDirectory.bin-runtime
var FeeCurrencyDirectoryBytecodeRaw []byte

//go:embed compiled/MockOracle.bin-runtime
var MockOracleBytecodeRaw []byte
