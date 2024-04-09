#!/bin/bash
#shellcheck disable=SC2034,SC2155,SC2086
set -xeo pipefail

export FEE_CURRENCY=$(\
	forge create --root . --contracts . --private-key $ACC_PRIVKEY DebugFeeCurrency.sol:DebugFeeCurrency --constructor-args '100000000000000000000000000' $1 $2 --json\
	| jq .deployedTo -r)

cast send --private-key $ACC_PRIVKEY $ORACLE3 'setExchangeRate(address, uint256, uint256)' $FEE_CURRENCY 2ether 1ether
cast send --private-key $ACC_PRIVKEY $FEE_CURRENCY_DIRECTORY_ADDR 'setCurrencyConfig(address, address, uint256)' $FEE_CURRENCY $ORACLE3 60000
echo Fee currency: $FEE_CURRENCY

(cd ../js-tests/ && ./send_tx.mjs "$(cast chain-id)" $ACC_PRIVKEY $FEE_CURRENCY)
