/*
 * assign stake reward
 */
const { BN, Long, bytes, units } = require('@zilliqa-js/util');
const { Zilliqa } = require('@zilliqa-js/zilliqa');
const { toBech32Address, getAddressFromPrivateKey } = require('@zilliqa-js/crypto');

const zilliqa = new Zilliqa('https://dev-api.zilliqa.com');
const CHAIN_ID = 2;
const MSG_VERSION = 1;
const VERSION = bytes.pack(CHAIN_ID, MSG_VERSION);

const PRIVATE_KEY = '';
const GAS_PRICE = units.toQa('1000', units.Units.Li);

const STAKING_PROXY_ADDR = toBech32Address("0123456789012345678901234567890123456789");


async function main() {
    zilliqa.wallet.addByPrivateKey(PRIVATE_KEY);
    const address = getAddressFromPrivateKey(PRIVATE_KEY);

    console.log("Invoking assign stake rewards...");
    console.log("Your account address is:");
    console.log(`${address}`);

    try {
        const contract = zilliqa.contracts.at(STAKING_PROXY_ADDR);
        const ssnList = [
            {
                "constructor": "SsnRewardShare",
                "argtypes": [],
                "arguments": [
                    "0x1234567890123456789012345678901234567890",
                    "50000000"
                ]
            }
        ];
        const callTx = await contract.call(
            'assign_stake_reward',
            [
                {
                    vname: 'ssnreward_list',
                    type: 'List SsnRewardShare',
                    value: ssnList
                },
                {
                    vname: 'reward_blocknum',
                    type: 'Uint32',
                    value: '50000'
                },
            ],
            {
                version: VERSION,
                amount: new BN(0), // sending amounts in ZIL, converted to Qa
                gasPrice: GAS_PRICE,
                gasLimit: Long.fromNumber(10000)
            },
            33,
            1000,
            true
        );
        console.log(JSON.stringify(callTx.receipt, null, 4));

    } catch (err) {
        console.log(err);
    }
}

main();