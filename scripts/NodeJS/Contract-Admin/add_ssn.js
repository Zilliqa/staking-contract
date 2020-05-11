/*
 * add ssn
 * used by admin only
 */
const { BN, Long, bytes, units } = require('@zilliqa-js/util');
const { Zilliqa } = require('@zilliqa-js/zilliqa');
const { toBech32Address, getAddressFromPrivateKey } = require('@zilliqa-js/crypto');

// change the following parameters
const API = 'http://localhost:5555'
const URL_API = API;
const URL_RAW = API;
const CHAIN_ID = 1;
const PRIVATE_KEY = 'e53d1c3edaffc7a7bab5418eb836cf75819a82872b4a1a0f1c7fcf5c3e020b89'; // admin
const STAKING_PROXY_ADDR = toBech32Address("0x35C36cEC66a7f5f5393f8b84eB56F4bd552dDb87"); // checksum proxy address
const SSN_ADDR = "0xf6dad9e193fa2959a849b81caf9cb6ecde466771" // ssn address to be registered with '0x'

const zilliqa = new Zilliqa(API);
const MSG_VERSION = 1;
const VERSION = bytes.pack(CHAIN_ID, MSG_VERSION);
const GAS_PRICE = units.toQa('1000', units.Units.Li);

async function main() {
    zilliqa.wallet.addByPrivateKey(PRIVATE_KEY);
    const address = getAddressFromPrivateKey(PRIVATE_KEY);
    console.log("Your account address is: %o", `${address}`);
    console.log("proxy: %o\n", STAKING_PROXY_ADDR);

    console.log("------------------------ begin add ssn ------------------------\n");
    try {
        const contract = zilliqa.contracts.at(STAKING_PROXY_ADDR);
        const callTx = await contract.call(
            'add_ssn',
            [
                {
                    vname: 'ssnaddr',
                    type: 'ByStr20',
                    value: `${SSN_ADDR}`
                },
                {
                    vname: 'urlraw',
                    type: 'String',
                    value: `${URL_RAW}`
                },
                {
                    vname: 'urlapi',
                    type: 'String',
                    value: `${URL_API}`
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
        console.log("------------------------ end add ssn ------------------------\n");

    } catch (err) {
        console.log(err);
    }
}

main();