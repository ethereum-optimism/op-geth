#!/bin/bash
set -eo pipefail

SCRIPT_DIR=$(readlink -f "$(dirname "$0")")

## Start geth
cd "$SCRIPT_DIR/.." || exit 1
make geth
trap 'kill %%' EXIT  # kill bg job at exit
build/bin/geth --dev --http --http.api eth,web3,net &> "$SCRIPT_DIR/geth.log" &

# Wait for geth to be ready
for _ in {1..10}
do
	if cast block &> /dev/null
	then
		break
	fi
	sleep 0.2
done

## Run tests
echo Geth ready, start tests
failures=0
tests=0
cd "$SCRIPT_DIR" || exit 1
for f in test_*
do
	echo -e "\nRun $f"
	if "./$f"
	then
		tput setaf 2 || true
		echo "PASS $f"
	else
		tput setaf 1 || true
		echo "FAIL $f ❌"
		((failures++)) || true
	fi
	tput init || true
	((tests++)) || true
done

## Final summary
echo
if [[ $failures -eq 0 ]] 
then
	tput setaf 2 || true
	echo All tests succeeded!
else
	tput setaf 1 || true
	echo $failures/$tests failed.
fi
tput init || true
exit $failures
