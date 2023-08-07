#/bin/bash
set -eo pipefail

# Before execution this script, run geth with a command like
# make geth && build/bin/geth --dev --http --http.api eth,web3,net console --exec "eth.sendTransaction({from: eth.coinbase, to: '0x42cf1bbc38BaAA3c4898ce8790e21eD2738c6A4a', value: web3.toWei(50, 'ether')}); admin.sleep(20)"

export ETH_RPC_URL=http://127.0.0.1:8545
ACC_ADDR=0x42cf1bbc38BaAA3c4898ce8790e21eD2738c6A4a
ACC_PRIVKEY=0x2771aff413cac48d9f8c114fabddd9195a2129f3c2c436caa07e27bb7f58ead5

REGISTRY_ADDR=0x000000000000000000000000000000000000ce10

# Deploy and register token
TOKEN_ADDR=$(forge create --private-key $ACC_PRIVKEY contracts/GoldToken.sol:GoldToken --constructor-args false | awk '/Deployed to/ {print $3}')
cast send --private-key $ACC_PRIVKEY $REGISTRY_ADDR 'function setAddressFor(string calldata identifier, address addr) external' GoldToken $TOKEN_ADDR
echo Account address: $ACC_ADDR, token address: $TOKEN_ADDR

# Send token and check balance
cast balance 0x000000000000000000000000000000000000dEaD
cast send --private-key $ACC_PRIVKEY $TOKEN_ADDR 'function transfer(address to, uint256 value) external returns (bool)' 0x000000000000000000000000000000000000dEaD 100
cast balance 0x000000000000000000000000000000000000dEaD
