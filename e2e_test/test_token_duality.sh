#!/bin/bash
#shellcheck disable=SC2086
set -eo pipefail

source shared.sh

# Register token
cast send --private-key $ACC_PRIVKEY $REGISTRY_ADDR 'function setAddressFor(string calldata identifier, address addr) external' GoldToken $TOKEN_ADDR

# Send token and check balance
balance_before=$(cast balance 0x000000000000000000000000000000000000dEaD)
cast send --private-key $ACC_PRIVKEY $TOKEN_ADDR 'function transfer(address to, uint256 value) external returns (bool)' 0x000000000000000000000000000000000000dEaD 100
balance_after=$(cast balance 0x000000000000000000000000000000000000dEaD)
echo "Balance change: $balance_before -> $balance_after"
[[ $((balance_before + 100)) -eq $balance_after ]] || (echo "Balance did not change as expected"; exit 1)
