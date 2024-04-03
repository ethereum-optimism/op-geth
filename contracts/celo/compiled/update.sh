#!/bin/bash

SCRIPT_DIR=$(readlink -f "$(dirname "$0")")
CONTRACTS_DIR=${CELO_OPTIMISM_REPO:-~/optimism}/packages/contracts-bedrock

(cd "$CONTRACTS_DIR" && forge build)

for contract in GoldToken CeloRegistry SortedOracles FeeCurrencyWhitelist FeeCurrency Proxy MockSortedOracles
do
	contract_json="$CONTRACTS_DIR/forge-artifacts/$contract.sol/$contract.json"
	jq .abi "$contract_json" > "$SCRIPT_DIR/$contract.abi"
	jq .deployedBytecode.object -r "$contract_json" > "$SCRIPT_DIR/$contract.bin-runtime"
done
