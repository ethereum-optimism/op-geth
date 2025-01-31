#!/usr/bin/env bash

# This script is used to sync superchain configs in the registry with OP Geth.

set -euo pipefail

# Constants
REGISTRY_COMMIT=$(cat superchain-registry-commit.txt)
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

repodir=$(mktemp -d)
workdir=$(mktemp -d)

# Clone the registry
echo "Cloning SR..."
cd "$repodir"
git clone --no-checkout --depth 1 --shallow-submodules https://github.com/ethereum-optimism/superchain-registry.git
cd "$repodir/superchain-registry"
git fetch --depth 1 origin "$REGISTRY_COMMIT"
git checkout "$REGISTRY_COMMIT"

echo "Copying configs..."
cp -r superchain/configs "$workdir/configs"
cp -r superchain/extra/genesis "$workdir/genesis"
cp -r superchain/extra/dictionary "$workdir/dictionary"

cd "$workdir"
echo "Using $workdir as workdir..."

# Create a simple mapping of chain id -> config name to make looking up chains by their ID easier.
echo "Generating index of configs..."

echo "{}" > chains.json

# Function to process each network directory
process_network_dir() {
    local network_dir="$1"
    local network_name=$(basename "$network_dir")

    echo "Processing chains in $network_name superchain..."

    # Find all TOML files in the network directory
    find "$network_dir" -type f -name "*.toml" | LC_ALL=C sort | while read -r toml_file; do
        if [[ "$toml_file" == "configs/$network_name/superchain.toml" ]]; then
            continue
        fi

        echo "Processing $toml_file..."
        # Extract chain_id from TOML file using dasel
        chain_id=$(dasel -f "$toml_file" -r toml "chain_id" | tr -d '"')

        # Skip if chain_id is empty
        if [ -z "$chain_id" ]; then
            continue
        fi

        # Create JSON object for this config
        config_json=$(jq -n \
            --arg name "$(basename "${toml_file%.*}")" \
            --arg network "$network_name" \
            '{
                "name": $name,
                "network": $network
            }')

        # Add this config to the result JSON using the chain_id as the key
        jq --argjson config "$config_json" \
           --arg chain_id "$chain_id" \
           '. + {($chain_id): $config}' chains.json > temp.json
        mv temp.json chains.json
    done
}

# Process each network directory in configs
for network_dir in configs/*; do
    if [ -d "$network_dir" ]; then
        process_network_dir "$network_dir"
    fi
done

# Archive the genesis configs as a ZIP file. ZIP is used since it can be efficiently used as a filesystem.
echo "Archiving configs..."
echo "$REGISTRY_COMMIT" > COMMIT
# We need to normalize the lastmod dates and permissions to ensure the ZIP file is deterministic.
find . -exec touch -t 198001010000.00 {} +
chmod -R 755 ./*
files=$(find . -type f | LC_ALL=C sort)
echo -n "$files" | xargs zip -9 -oX --quiet superchain-configs.zip
zipinfo superchain-configs.zip
mv superchain-configs.zip "$SCRIPT_DIR/superchain/superchain-configs.zip"

echo "Cleaning up..."
rm -rf "$repodir"
rm -rf "$workdir"

echo "Done."