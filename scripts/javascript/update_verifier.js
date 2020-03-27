/*
 * update verifier
 */
const { BN, Long, bytes, units } = require('@zilliqa-js/util');
const { Zilliqa } = require('@zilliqa-js/zilliqa');
const { toBech32Address, getAddressFromPrivateKey } = require('@zilliqa-js/crypto');

// change the following parameters
const API = 'http://localhost:5555'
const CHAIN_ID = 1;
const PRIVATE_KEY = 'e53d1c3edaffc7a7bab5418eb836cf75819a82872b4a1a0f1c7fcf5c3e020b89'; // admin
const STAKING_PROXY_ADDR = toBech32Address("0xDB5Dc7118765A84B6c6A582280fA37A1DD2d9f69"); // checksum proxy address
const NEW_VERIFIER = '0x381f4008505e940ad7681ec3468a719060caf796'; // new verifier checksum address with '0x'

const zilliqa = new Zilliqa(API);
const MSG_VERSION = 1;
const VERSION = bytes.pack(CHAIN_ID, MSG_VERSION);
const GAS_PRICE = units.toQa('1000', units.Units.Li);


async function main() {
    zilliqa.wallet.addByPrivateKey(PRIVATE_KEY);
    const address = getAddressFromPrivateKey(PRIVATE_KEY);
    console.log("Your account address is: %o", `${address}`);
    console.log("proxy: %o\n", STAKING_PROXY_ADDR);

    console.log("------------------------ begin update verifier ------------------------\n");
    try {
        const contract = zilliqa.contracts.at(STAKING_PROXY_ADDR);
        const callTx = await contract.call(
            'update_verifier',
            [
                {
                    vname: 'verif',
                    type: 'ByStr20',
                    value: `${NEW_VERIFIER}`
                }
            ],
            {
                version: VERSION,
                amount: new BN(0),
                gasPrice: GAS_PRICE,
                gasLimit: Long.fromNumber(10000)
            },
            33,
            1000,
            true
        );
        console.log("transaction: %o", callTx.id);
        console.log(JSON.stringify(callTx.receipt, null, 4));

    } catch (err) {
        console.log(err);
    }
    console.log("------------------------ end update verifier ------------------------\n");
}

main();