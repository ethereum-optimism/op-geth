#!/usr/bin/env node
import { createPublicClient, createWalletClient, http, defineChain } from 'viem'
import { celoAlfajores } from 'viem/chains'
import { privateKeyToAccount } from 'viem/accounts'

const [chainId, privateKey, feeCurrency] = process.argv.slice(2)
const devChain = defineChain({
  ...celoAlfajores,
  id: parseInt(chainId, 10),
  name: 'local dev chain',
  network: 'dev',
  rpcUrls: {
    default: {
      http: ['http://127.0.0.1:8545'],
    },
  },
})

const account = privateKeyToAccount(privateKey) 
const walletClient = createWalletClient({
  account,
  chain: devChain,
  transport: http(),
})

const request = await walletClient.prepareTransactionRequest({
  account,
  to: '0x00000000000000000000000000000000DeaDBeef',
  value: 2,
  gas: 90000,
  feeCurrency,
  maxFeePerGas: 2000000000n,
  maxPriorityFeePerGas: 0n,
})
const signature = await walletClient.signTransaction(request)
const hash = await walletClient.sendRawTransaction({ serializedTransaction: signature }) 
console.log(hash)
