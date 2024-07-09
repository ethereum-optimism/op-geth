#!/bin/bash

SCRIPT_DIR=$(readlink -f "$(dirname "$0")")

CONTRACTS_DIR=${CELO_OPTIMISM_REPO:-~/optimism}/packages/contracts-bedrock
forge build --root "$CONTRACTS_DIR"

for contract in FeeCurrency
do
	contract_json="$CONTRACTS_DIR/forge-artifacts/$contract.sol/$contract.json"
	jq .abi "$contract_json" > "$SCRIPT_DIR/$contract.abi"
	jq .deployedBytecode.object -r "$contract_json" > "$SCRIPT_DIR/$contract.bin-runtime"
done

CONTRACTS_DIR=${CELO_MONOREPO:-~/celo-monorepo}/packages/protocol
forge build --root "$CONTRACTS_DIR"

for contract in GoldToken FeeCurrencyDirectory IFeeCurrencyDirectory MockOracle
do
	contract_json="$CONTRACTS_DIR/out/$contract.sol/$contract.json"
	jq .abi "$contract_json" > "$SCRIPT_DIR/$contract.abi"
	jq .deployedBytecode.object -r "$contract_json" > "$SCRIPT_DIR/$contract.bin-runtime"
done

# We only need the abi for the interface (IFeeCurrencyDirectory) and the
# bytecode for the implementation (FeeCurrencyDirectory), so let's delete the other.
rm "$SCRIPT_DIR/IFeeCurrencyDirectory.bin-runtime" "$SCRIPT_DIR/FeeCurrencyDirectory.abi"
