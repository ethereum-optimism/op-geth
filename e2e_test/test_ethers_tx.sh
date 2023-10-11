#!/bin/bash
set -eo pipefail

source shared.sh
prepare_node

cd js-tests && ./node_modules/mocha/bin/mocha.js test_ethers_tx.mjs
