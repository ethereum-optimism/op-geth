package celo

import _ "embed"

//go:generate go run ../../cmd/abigen --pkg abigen --out abigen/FeeCurrency.go --abi compiled/FeeCurrency.abi --type FeeCurrency
//go:generate go run ../../cmd/abigen --pkg abigen --out abigen/FeeCurrencyDirectory.go --abi compiled/IFeeCurrencyDirectory.abi --type FeeCurrencyDirectory

//go:embed compiled/GoldToken.bin-runtime
var CeloTokenBytecodeRaw []byte

//go:embed compiled/FeeCurrency.bin-runtime
var FeeCurrencyBytecodeRaw []byte

//go:embed compiled/FeeCurrencyDirectory.bin-runtime
var FeeCurrencyDirectoryBytecodeRaw []byte

//go:embed compiled/MockOracle.bin-runtime
var MockOracleBytecodeRaw []byte
