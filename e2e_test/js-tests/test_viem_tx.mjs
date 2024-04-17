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

const testNonceBump = async (
	firstCap,
	firstCurrency,
	secondCap,
	secondCurrency,
	shouldReplace,
) => {
	const syncBarrierRequest = await walletClient.prepareTransactionRequest({
		account,
		to: "0x00000000000000000000000000000000DeaDBeef",
		value: 2,
		gas: 22000,
	});
	const firstTxHash = await walletClient.sendTransaction({
		account,
		to: "0x00000000000000000000000000000000DeaDBeef",
		value: 2,
		gas: 90000,
		maxFeePerGas: firstCap,
		maxPriorityFeePerGas: firstCap,
		nonce: syncBarrierRequest.nonce + 1,
		feeCurrency: firstCurrency,
	});
	var secondTxHash;
	try {
		secondTxHash = await walletClient.sendTransaction({
			account,
			to: "0x00000000000000000000000000000000DeaDBeef",
			value: 3,
			gas: 90000,
			maxFeePerGas: secondCap,
			maxPriorityFeePerGas: secondCap,
			nonce: syncBarrierRequest.nonce + 1,
			feeCurrency: secondCurrency,
		});
	} catch (err) {
		// If shouldReplace, no error should be thrown
		// If shouldReplace == false, exactly the underpriced error should be thrown
		if (
			err.cause.details != "replacement transaction underpriced" ||
			shouldReplace
		) {
			throw err; // Only throw if unexpected error.
		}
	}
	const syncBarrierSignature =
		await walletClient.signTransaction(syncBarrierRequest);
	const barrierTxHash = await walletClient.sendRawTransaction({
		serializedTransaction: syncBarrierSignature,
	});
	await publicClient.waitForTransactionReceipt({ hash: barrierTxHash });
	if (shouldReplace) {
		// The new transaction was included.
		await publicClient.waitForTransactionReceipt({ hash: secondTxHash });
	} else {
		// The original transaction was not replaced.
		await publicClient.waitForTransactionReceipt({ hash: firstTxHash });
	}
};

describe("viem send tx", () => {
	it("send basic tx and check receipt", async () => {
		const request = await walletClient.prepareTransactionRequest({
			account,
			to: "0x00000000000000000000000000000000DeaDBeef",
			value: 1,
			gas: 21000,
		});
		const signature = await walletClient.signTransaction(request);
		const hash = await walletClient.sendRawTransaction({
			serializedTransaction: signature,
		});
		const receipt = await publicClient.waitForTransactionReceipt({ hash });
		assert.equal(receipt.status, "success", "receipt status 'failure'");
	}).timeout(10_000);

	it("send tx with gas estimation and check receipt", async () => {
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
		assert.equal(receipt.status, "success", "receipt status 'failure'");
	}).timeout(10_000);

	it("send fee currency tx and check receipt", async () => {
		const request = await walletClient.prepareTransactionRequest({
			account,
			to: "0x00000000000000000000000000000000DeaDBeef",
			value: 2,
			gas: 90000,
			feeCurrency: process.env.FEE_CURRENCY,
			maxFeePerGas: 2000000000n,
			maxPriorityFeePerGas: 0n,
		});
		const signature = await walletClient.signTransaction(request);
		const hash = await walletClient.sendRawTransaction({
			serializedTransaction: signature,
		});
		const receipt = await publicClient.waitForTransactionReceipt({ hash });
		assert.equal(receipt.status, "success", "receipt status 'failure'");
	}).timeout(10_000);

	it("test gas price difference for fee currency", async () => {
		const request = await walletClient.prepareTransactionRequest({
			account,
			to: "0x00000000000000000000000000000000DeaDBeef",
			value: 2,
			gas: 90000,
			feeCurrency: process.env.FEE_CURRENCY,
		});

		const gasPriceNative = await publicClient.getGasPrice({});
		var maxPriorityFeePerGasNative =
			await publicClient.estimateMaxPriorityFeePerGas({});
		const block = await publicClient.getBlock({});
		assert.equal(
			BigInt(block.baseFeePerGas) + maxPriorityFeePerGasNative,
			gasPriceNative,
		);

		// viem's getGasPrice does not expose additional request parameters,
		// but Celo's override 'chain.fees.estimateFeesPerGas' action does.
		// this will call the eth_gasPrice and eth_maxPriorityFeePerGas methods
		// with the additional feeCurrency parameter internally
		var fees = await publicClient.estimateFeesPerGas({
			type: "eip1559",
			request: {
				feeCurrency: process.env.FEE_CURRENCY,
			},
		});
		// first check that the fee currency denominated gas price
		// converts properly to the native gas price
		assert.equal(fees.maxFeePerGas, gasPriceNative * 2n);
		assert.equal(fees.maxPriorityFeePerGas, maxPriorityFeePerGasNative * 2n);

		// check that the prepared transaction request uses the
		// converted gas price internally
		assert.equal(request.maxFeePerGas, fees.maxFeePerGas);
		assert.equal(request.maxPriorityFeePerGas, fees.maxPriorityFeePerGas);
	}).timeout(10_000);

	it("send overlapping nonce tx in different currencies", async () => {
		const priceBump = 1.1;
		const rate = 2;
		// Native to FEE_CURRENCY
		const nativeCap = 30_000_000_000;
		const bumpCurrencyCap = BigInt(Math.round(nativeCap * rate * priceBump));
		const failToBumpCurrencyCap = BigInt(
			Math.round(nativeCap * rate * priceBump) - 1,
		);
		// FEE_CURRENCY to Native
		const currencyCap = 60_000_000_000;
		const bumpNativeCap = BigInt(Math.round((currencyCap * priceBump) / rate));
		const failToBumpNativeCap = BigInt(
			Math.round((currencyCap * priceBump) / rate) - 1,
		);
		const tokenCurrency = process.env.FEE_CURRENCY;
		const nativeCurrency = null;
		await testNonceBump(
			nativeCap,
			nativeCurrency,
			bumpCurrencyCap,
			tokenCurrency,
			true,
		);
		await testNonceBump(
			nativeCap,
			nativeCurrency,
			failToBumpCurrencyCap,
			tokenCurrency,
			false,
		);
		await testNonceBump(
			currencyCap,
			tokenCurrency,
			bumpNativeCap,
			nativeCurrency,
			true,
		);
		await testNonceBump(
			currencyCap,
			tokenCurrency,
			failToBumpNativeCap,
			nativeCurrency,
			false,
		);
	}).timeout(10_000);
});
