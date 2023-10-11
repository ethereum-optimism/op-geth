import { assert } from "chai";
import "mocha";
import {
	createPublicClient,
	createWalletClient,
	http,
	defineChain,
} from "viem";
import { celoAlfajores } from "viem/chains";
import { privateKeyToAccount } from "viem/accounts";

// Setup up chain
const devChain = defineChain({
	...celoAlfajores,
	id: 1337,
	name: "local dev chain",
	network: "dev",
	rpcUrls: {
		default: {
			http: [process.env.ETH_RPC_URL],
		},
	},
});

// Set up clients/wallet
const publicClient = createPublicClient({
	chain: devChain,
	transport: http(),
});
const account = privateKeyToAccount(process.env.ACC_PRIVKEY);
const walletClient = createWalletClient({
	account,
	chain: devChain,
	transport: http(),
});

describe("viem send tx", () => {
	it("send basic tx and check receipt", async () => {
		const request = await walletClient.prepareTransactionRequest({
			account,
			to: "0x00000000000000000000000000000000DeaDBeef",
			value: 1,
		});
		const signature = await walletClient.signTransaction(request);
		const hash = await walletClient.sendRawTransaction({
			serializedTransaction: signature,
		});
		const receipt = await publicClient.waitForTransactionReceipt({ hash });
	}).timeout(10_000);

	it("send fee currency tx and check receipt", async () => {
		const request = await walletClient.prepareTransactionRequest({
			account,
			to: "0x00000000000000000000000000000000DeaDBeef",
			value: 2,
			gas: 90000,
			feeCurrency: process.env.FEE_CURRENCY,
		});
		const signature = await walletClient.signTransaction(request);
		const hash = await walletClient.sendRawTransaction({
			serializedTransaction: signature,
		});
		const receipt = await publicClient.waitForTransactionReceipt({ hash });
	}).timeout(10_000);
});
