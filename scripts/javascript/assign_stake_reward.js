/*
 * assign stake reward
 * used by verifier only
 */
const { BN, Long, bytes, units } = require('@zilliqa-js/util');
const { Zilliqa } = require('@zilliqa-js/zilliqa');
const { toBech32Address, getAddressFromPrivateKey } = require('@zilliqa-js/crypto');

// change the following parameters
const API = 'http://localhost:5555'
const CHAIN_ID = 1;
const PRIVATE_KEY = 'd96e9eb5b782a80ea153c937fa83e5948485fbfc8b7e7c069d7b914dbc350aba'; // verifier
const STAKING_PROXY_ADDR = toBech32Address("0xDB5Dc7118765A84B6c6A582280fA37A1DD2d9f69"); // checksum proxy address
const REWARD_BLOCKNUM = 50000; // tx block num when ssns were verified
// in the arguments, the second parameter '50' is the percentage of stake_deposit to reward
const SSN_REWARD_LIST = [
    {
        "constructor": "SsnRewardShare",
        "argtypes": "[]",
        "arguments": [
            "0xf6dad9e193fa2959a849b81caf9cb6ecde466771",
            "50"
        ]
    }
];

const zilliqa = new Zilliqa(API);
const MSG_VERSION = 1;
const VERSION = bytes.pack(CHAIN_ID, MSG_VERSION);
const GAS_PRICE = units.toQa('1000', units.Units.Li);


async function main() {
    zilliqa.wallet.addByPrivateKey(PRIVATE_KEY);
    const address = getAddressFromPrivateKey(PRIVATE_KEY);
    console.log("Your account address is: %o", `${address}`);
    console.log("proxy: %o\n", STAKING_PROXY_ADDR);

    console.log("------------------------ begin assign stake reward ------------------------\n");
    try {
        const contract = zilliqa.contracts.at(STAKING_PROXY_ADDR);
        const callTx = await contract.call(
            'assign_stake_reward',
            [
                {
                    vname: 'ssnreward_list',
                    type: 'List SsnRewardShare',
                    value: SSN_REWARD_LIST
                },
                {
                    vname: 'reward_blocknum',
                    type: 'Uint32',
                    value: `${REWARD_BLOCKNUM}`
                },
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
    console.log("------------------------ end assign stake reward ------------------------\n");
}

main();