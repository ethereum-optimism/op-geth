package contracts

import _ "embed"

//go:generate go run ../../cmd/abigen --pkg abigen --out abigen/GoldToken.go --abi compiled/GoldToken.abi --type GoldToken
//go:generate go run ../../cmd/abigen --pkg abigen --out abigen/Registry.go --abi compiled/Registry.abi --type Registry

//go:embed compiled/Registry.bin-runtime
var RegistryBytecodeRaw []byte

//go:embed compiled/Proxy.bin-runtime
var ProxyBytecodeRaw []byte
