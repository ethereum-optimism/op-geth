#!/bin/bash
#shellcheck disable=SC2086
set -eo pipefail

source shared.sh

# Expect that the debitGasFees fails during tx submission
(cd debug-fee-currency && ./deploy_and_send_tx.sh true false) &> debug-fee-currency/send_tx.log || true
grep "debitGasFees reverted: This DebugFeeCurrency always fails in debitGasFees!" debug-fee-currency/send_tx.log \
  || (cat debug-fee-currency/send_tx.log && false)
