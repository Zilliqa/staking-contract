/*
 * update contract max stake
 * used by admin only
 */
const { BN, Long, bytes, units } = require('@zilliqa-js/util');
const { Zilliqa } = require('@zilliqa-js/zilliqa');
const { toBech32Address, getAddressFromPrivateKey } = require('@zilliqa-js/crypto');

// change the following parameters
const API = 'http://localhost:5555'
const CHAIN_ID = 1;
const PRIVATE_KEY = 'e53d1c3edaffc7a7bab5418eb836cf75819a82872b4a1a0f1c7fcf5c3e020b89'; // admin
const STAKING_PROXY_ADDR = toBech32Address("0x651b97542A0B339052d61eB13f6c4FcDBA1a0172"); // checksum proxy address
const CONTRACT_MAX_STAKE = units.toQa('4000', units.Units.Zil); // contract max stake amount in ZIL converted to Qa
const MAX_STAKE = units.toQa('3000', units.Units.Zil); // max stake amount in ZIL converted to Qa
const MIN_STAKE = units.toQa('10', units.Units.Zil); // min stake amount in ZIL converted to Qa

const zilliqa = new Zilliqa(API);
const MSG_VERSION = 1;
const VERSION = bytes.pack(CHAIN_ID, MSG_VERSION);
const GAS_PRICE = units.toQa('1000', units.Units.Li);


async function main() {
    zilliqa.wallet.addByPrivateKey(PRIVATE_KEY);
    const address = getAddressFromPrivateKey(PRIVATE_KEY);
    console.log("Your account address is: %o", `${address}`);
    console.log("proxy: %o\n", STAKING_PROXY_ADDR);

    console.log("------------------------ begin update staking parameter ------------------------\n");
    try {
        const contract = zilliqa.contracts.at(STAKING_PROXY_ADDR);
        const callTx = await contract.call(
            'update_staking_parameter',
            [
                {
                    vname: 'min_stake',
                    type: 'Uint128',
                    value: `${MIN_STAKE}`
                },
                {
                    vname: 'max_stake',
                    type: 'Uint128',
                    value: `${MAX_STAKE}`
                },
                {
                    vname: 'contract_max_stake',
                    type: 'Uint128',
                    value: `${CONTRACT_MAX_STAKE}`
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
    console.log("------------------------ end update staking parameter ------------------------\n");
}

main();