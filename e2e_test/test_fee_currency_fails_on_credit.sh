#!/bin/bash
#shellcheck disable=SC2086
set -eo pipefail

source shared.sh

# Expect that the creditGasFees failed and is logged by geth
tail -f -n0 geth.log >debug-fee-currency/geth.partial.log & # start log capture
(cd debug-fee-currency && ./deploy_and_send_tx.sh false true false)
sleep 0.5
kill %1 # stop log capture
grep "This DebugFeeCurrency always fails in (old) creditGasFees!" debug-fee-currency/geth.partial.log
