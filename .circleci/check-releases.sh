#!/bin/bash
set -euo pipefail

LATEST_RELEASE=$(curl -s --fail -L \
  -H "Accept: application/vnd.github+json" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/ethereum/go-ethereum/releases \
  | jq -r '(.[] | select(.draft==false) | select(.prerelease==false)).tag_name' | head -n 1)

echo "Detected latest go-ethereum release as ${LATEST_RELEASE}"

git remote add upstream https://github.com/ethereum/go-ethereum
git fetch upstream > /dev/null

if git branch --contains "${LATEST_RELEASE}" 2>&1 | grep -e '^[ *]*optimism$' > /dev/null
then
  echo "Up to date with latest release.  Great job! 🎉"
else
  echo "Release has not been merged"
  exit 1
fi
