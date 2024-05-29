package addresses

import "github.com/ethereum/go-ethereum/common"

var (
	RegistryAddress                = common.HexToAddress("0x000000000000000000000000000000000000ce10")
	GoldTokenAddress               = common.HexToAddress("0x471ece3750da237f93b8e339c536989b8978a438")
	FeeHandlerAddress              = common.HexToAddress("0xcd437749e43a154c07f3553504c68fbfd56b8778")
	MentoFeeHandlerSellerAddress   = common.HexToAddress("0x4efa274b7e33476c961065000d58ee09f7921a74")
	UniswapFeeHandlerSellerAddress = common.HexToAddress("0xd3aee28548dbb65df03981f0dc0713bfcbd10a97")
	FeeCurrencyDirectoryAddress    = common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb") // tmp address, real one no known yet
)
