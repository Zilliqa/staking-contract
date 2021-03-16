/*
 * get the amount deposited in the stake
 */
const { Zilliqa } = require('@zilliqa-js/zilliqa');
const { BN, Long, bytes, units } = require('@zilliqa-js/util');
const { toBech32Address, getAddressFromPrivateKey } = require('@zilliqa-js/crypto');

// change the following parameters
const API = 'https://dev-api.zilliqa.com'
const PRIVATE_KEY = 'e19d05c5452598e24caad4a0d85a49146f7be089515c905ae6a19e8a578a6930'; // private key of staker
const STAKING_PROXY_ADDR = toBech32Address("0x05d7e121e205a84bf1da2d60ac8a2484800fffb3"); // checksum proxy address

const zilliqa = new Zilliqa(API);
const CHAIN_ID = 333;
const MSG_VERSION = 1;
const VERSION = bytes.pack(CHAIN_ID, MSG_VERSION);
const GAS_PRICE = units.toQa('2000', units.Units.Li);

async function main() {
    zilliqa.wallet.addByPrivateKey(PRIVATE_KEY);
    const address = getAddressFromPrivateKey(PRIVATE_KEY).toLowerCase();
    console.log("Your account address is: %o", `${address}`);
    console.log("proxy: %o\n", STAKING_PROXY_ADDR);

    console.log("------------------------ begin delegate stake ------------------------\n");
    try {
        const contract = zilliqa.contracts.at(STAKING_PROXY_ADDR);
        const callTx = await contract.callWithoutConfirm(
            'ReDelegateStake',
            [
                {
                    vname: 'ssnaddr',
                    type: 'ByStr20',
                    value: '0x1b825f2bfe4515d34262c85e609313d8f88b2ae9'
                },
                {
                    vname: 'to_ssn',
                    type: 'ByStr20',
                    value: '0x227227e7eb01bf4fe00817fb6ed948d8bbcb4ad7'
                },
                {
                    vname: 'amount',
                    type: 'Uint128',
                    value: '10000000000000'
                }
            ],
            {
                version: VERSION,
                amount: new BN(0),
                gasPrice: GAS_PRICE,
                gasLimit: Long.fromNumber(30000)
            },
            false,
        );
        // process confirm
        console.log(`The transaction id is:`, callTx.id);
        console.log(`Waiting transaction be confirmed`);
        const confirmedTxn = await callTx.confirm(callTx.id);
        console.log(`The transaction status is:`);
        console.log(confirmedTxn.receipt);
    } catch (err) {
        console.log(err);
    }
    console.log("------------------------ end delegate stake ------------------------\n");
}

main();