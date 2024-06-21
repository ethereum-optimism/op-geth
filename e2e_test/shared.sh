#!/bin/bash
#shellcheck disable=SC2034  # unused vars make sense in a shared file

export ETH_RPC_URL=http://127.0.0.1:8545

export ACC_ADDR=0x42cf1bbc38BaAA3c4898ce8790e21eD2738c6A4a
export ACC_PRIVKEY=0x2771aff413cac48d9f8c114fabddd9195a2129f3c2c436caa07e27bb7f58ead5
export REGISTRY_ADDR=0x000000000000000000000000000000000000ce10
export TOKEN_ADDR=0x471ece3750da237f93b8e339c536989b8978a438
export FEE_CURRENCY_DIRECTORY_ADDR=0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
export FEE_CURRENCY=0x000000000000000000000000000000000000ce16
export FEE_HANDLER=0xcd437749e43a154c07f3553504c68fbfd56b8778
export ORACLE=0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb0001

export FIXIDITY_1=1000000000000000000000000
export ZERO_ADDRESS=0x0000000000000000000000000000000000000000

prepare_node () {
  (cd js-tests || exit 1; [[ -d node_modules ]] || npm install)
}