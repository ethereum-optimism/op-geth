import { assert } from "chai";
import "mocha";
import { ethers } from "ethers";

const provider = new ethers.JsonRpcProvider(process.env.ETH_RPC_URL);
const signer = new ethers.Wallet(process.env.ACC_PRIVKEY, provider);

describe("ethers.js send tx", () => {
	it("send basic tx and check receipt", async () => {
		const tx = await signer.sendTransaction({
			to: "0x00000000000000000000000000000000DeaDBeef",
			value: 1,
		});
		const receipt = await tx.wait();
	}).timeout(10_000);
});

describe("ethers.js compatibility tests with state", () => {
	it("provider.getBlock works (block has gasLimit set)", async () => {
		let block = await provider.getBlock();

		// These assertions trigger on undefined or null
		assert.notEqual(block, null);
		assert.notEqual(block.gasLimit, null);
	});

	it("EIP-1559 transactions supported (can get feeData)", async () => {
		const feeData = await provider.getFeeData();

		// These assertions trigger on undefined or null
		assert.notEqual(feeData, null);
		assert.notEqual(feeData.maxFeePerGas, null);
		assert.notEqual(feeData.maxPriorityFeePerGas, null);
		assert.notEqual(feeData.gasPrice, null);
	});

	it("block has gasLimit", async () => {
		const fullBlock = await provider.send("eth_getBlockByNumber", [
			"latest",
			true,
		]);
		assert.isTrue(fullBlock.hasOwnProperty("gasLimit"));
	});

	it("block has baseFeePerGas", async () => {
		const fullBlock = await provider.send("eth_getBlockByNumber", [
			"latest",
			true,
		]);
		assert.isTrue(fullBlock.hasOwnProperty("baseFeePerGas"));
	});
});
