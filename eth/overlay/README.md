# Overlays
Overlays allow you to add your custom logic to already deployed contracts and simulate events and calls on top of them.
With overlays you can create new view functions, modify existing ones, change field visibility, emit new events and query the historical data of many contracts with your modified source code.

## API
This explains how to use the overlay API.

### `overlay_callConstructor`
This method needs to be called once with the new bytecode.

It first does a lookup of the creationTx for the given contract.
Once it's found, it injects the new code and returns the new creation bytecode result from the EVM to the caller.

Example request:
```json
{
  "id" : "1",
  "jsonrpc" : "2.0",
  "method" : "overlay_callConstructor",
  "params" : [
    "<CONTRACT_ADDRESS>",
    "<BYTECODE>"
  ]
}
```

Example response:
```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "result": {
    "code": "<CREATION_BYTECODE>"
  }
}
```

### `overlay_getLogs`
This method can be called multiple times to receive new logs from your new bytecode.

It has the same interface as `eth_getLogs` but it also accepts state overrides as the second param.
We can pass the creation bytecode from the call to `overlay_callConstructor` along to `overlay_getLogs` as state overrides.
The passed block range for the filter defines the initial block range that needs to be replayed with the given state overrides.
Once all blocks are replayed, the logs are returned to the caller.

Example request:
```json
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "method" : "overlay_getLogs",
   "params" : [
      {
         "address" : "<CONTRACT_ADDRESS>",
         "fromBlock" : "0x6e7dd00",
         "toBlock" : "0x6e7dd00"
      },
      {
         "<CONTRACT_ADDRESS>" : {
            "code" : "<CREATION_BYTECODE>"
        }
      }
    ]
}
```

Example response as in `eth_getLogs`:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": [
    {
      "address": "<CONTRACT_ADDRESS",
      "topics": [
        "0xaabd75b90fb7114eb9587a54f00ce5ebe8cb4a70627f3a6c26e506ffd771fe2f",
        "0x00000000000000000000000000000000fc25870c6ed6b6c7e41fb078b7656f69",
        "0x000000000000000000000000000000000000000000000000000000000004d77c"
      ],
      "data": "0x0000000000000000000000000000000000000000000000000000000000000001",
      "blockNumber": "0x6e7dd00",
      "transactionHash": "0x32983cb59ae25889efac5ec4850534421495b1a231c8472e822a340eff8db23e",
      "transactionIndex": "0x1",
      "blockHash": "0xffe6aa3ba2bcb08f8fdd340e0af309e54d2da0b7d08efb3185f6dece00d0a3c6",
      "logIndex": "0x1",
      "removed": false
    }
  ]
}

```

### `eth_call`
This method can be called multiple times to call new view functions that you defined in your new bytecode.

By sending the creation bytecode received from `overlay_callConstructor` as state overrides to `eth_call` you'll be able to call new functions on your contract.

## Tests
There's a [postman collection for overlays](Overlay_Testing.json) with sample requests for `overlay_callConstructor` and `overlay_getLogs` which can be used for reference and refactorings.

## Configuration
- add `overlay` to your `--http.api` flag
