/*
 * pause
 */
const { BN, Long, bytes, units } = require('@zilliqa-js/util');
const { Zilliqa } = require('@zilliqa-js/zilliqa');
const { toBech32Address, getAddressFromPrivateKey } = require('@zilliqa-js/crypto');

const zilliqa = new Zilliqa('https://dev-api.zilliqa.com');
const zilliqa = new Zilliqa(API);
const CHAIN_ID = 2;
const MSG_VERSION = 1;
const VERSION = bytes.pack(CHAIN_ID, MSG_VERSION);

const PRIVATE_KEY = '';
const GAS_PRICE = units.toQa('1000', units.Units.Li);

const STAKING_PROXY_ADDR = toBech32Address("0123456789012345678901234567890123456789");

async function main() {
    zilliqa.wallet.addByPrivateKey(PRIVATE_KEY);
    const address = getAddressFromPrivateKey(PRIVATE_KEY);

    console.log("Invoking pause at: %o", STAKING_PROXY_ADDR);
    console.log("Network: %o", API);
    console.log("Private key: %o", PRIVATE_KEY);
    console.log("Account address: %o", `${address}`);

    try {
        const balance = await zilliqa.blockchain.getBalance(address);
        const nonce = balance.result;

        const contract = zilliqa.contracts.at(STAKING_PROXY_ADDR);
        const callTx = await contract.call(
            'unpause',
            [],
            {
                version: VERSION,
                amount: new BN(0),
                gasPrice: GAS_PRICE,
                gasLimit: Long.fromNumber(10000),
                nonce: (nonce + 1)
            },
            true
        );
        console.log(JSON.stringify(callTx.receipt, null, 4));
    } catch (err) {
        console.log(err);
    }
}

main();