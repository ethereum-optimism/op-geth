#!/bin/bash
set -eo pipefail

source shared.sh
prepare_node

# Add our account as oracle and submit value
cast send --private-key "$ACC_PRIVKEY" "$SORTED_ORACLES_ADDR" 'addOracle(address token, address oracleAddress)' "$FEE_CURRENCY" "$ACC_ADDR"
cast send --private-key "$ACC_PRIVKEY" "$SORTED_ORACLES_ADDR" 'report(address token, uint256 value, address lesserKey, address greaterKey)' "$FEE_CURRENCY" 2000000000000000000000000 "$ZERO_ADDRESS" "$ZERO_ADDRESS"

cd js-tests && ./node_modules/mocha/bin/mocha.js test_viem_tx.mjs
