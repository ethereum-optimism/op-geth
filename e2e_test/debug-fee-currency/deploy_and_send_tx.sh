#!/bin/bash
#shellcheck disable=SC2034,SC2155,SC2086
set -xeo pipefail

export FEE_CURRENCY=$(\
	forge create --root . --contracts . --private-key $ACC_PRIVKEY DebugFeeCurrency.sol:DebugFeeCurrency --constructor-args '100000000000000000000000000' $1 $2 --json\
	| jq .deployedTo -r)

cast send --private-key $ACC_PRIVKEY $FEE_CURRENCY_DIRECTORY_ADDR 'setCurrencyConfig(address, address, address, uint256)' $FEE_CURRENCY $FEE_CURRENCY $ORACLE 60000 --legacy
cast send --private-key $ACC_PRIVKEY $SORTED_ORACLES_ADDR 'setMedianRate(address, uint256)' $FEE_CURRENCY 2000000000000000000000000 --legacy
echo Fee currency: $FEE_CURRENCY

(cd ../js-tests/ && ./send_tx.mjs "$(cast chain-id)" $ACC_PRIVKEY $FEE_CURRENCY)
