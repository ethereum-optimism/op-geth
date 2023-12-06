#!/bin/bash
#shellcheck disable=SC2086
set -eo pipefail
set -x

source shared.sh

# Send token and check balance
balance_before=$(cast balance $FEE_HANDLER)
tx_json=$(cast send --json --private-key $ACC_PRIVKEY $TOKEN_ADDR 'transfer(address to, uint256 value) returns (bool)' 0x000000000000000000000000000000000000dEaD 100)
gas_used=$(echo $tx_json | jq -r '.gasUsed')
block_number=$(echo $tx_json | jq -r '.blockNumber')
base_fee=$(cast base-fee $block_number)
expected_balance_change=$((base_fee * gas_used))
balance_after=$(cast balance $FEE_HANDLER)
echo "Balance change: $balance_before -> $balance_after"
[[ $((balance_before + expected_balance_change)) -eq $balance_after ]] || (echo "Balance did not change as expected"; exit 1)