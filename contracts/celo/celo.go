package contracts

//go:generate go run ../../cmd/abigen --pkg abigen --out abigen/GoldToken.go --abi compiled/GoldToken.abi --type GoldToken
//go:generate go run ../../cmd/abigen --pkg abigen --out abigen/Registry.go --abi compiled/Registry.abi --type Registry
