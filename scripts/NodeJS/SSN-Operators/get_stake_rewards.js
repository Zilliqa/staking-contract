/*
 * get the stake rewards
 */
const { Zilliqa } = require('@zilliqa-js/zilliqa');
const { toBech32Address, getAddressFromPrivateKey } = require('@zilliqa-js/crypto');

// change the following parameters
const API = 'http://localhost:5555'
const PRIVATE_KEY = '589417286a3213dceb37f8f89bd164c3505a4cec9200c61f7c6db13a30a71b45'; // private key of staker
const STAKING_PROXY_ADDR = toBech32Address("0xF8AE69E4d0a4f073Bd2223a042ce89226E6d3663"); // checksum proxy address

const zilliqa = new Zilliqa(API);

async function main() {
    zilliqa.wallet.addByPrivateKey(PRIVATE_KEY);
    const address = getAddressFromPrivateKey(PRIVATE_KEY).toLowerCase();
    console.log("Your account address is: %o", `${address}`);
    console.log("proxy: %o\n", STAKING_PROXY_ADDR);

    console.log("------------------------ begin get stake rewards amount ------------------------\n");
    try {
        const contract = zilliqa.contracts.at(STAKING_PROXY_ADDR);
        const state = await contract.getState();
        const implAddress = state.implementation;
        
        const implContract = zilliqa.contracts.at(toBech32Address(implAddress));
        const implState = await implContract.getState();

        if (Object.keys(implState.ssnlist).length === 0) {
            console.error("SSN list is empty");
        } else if (!implState.ssnlist.hasOwnProperty(address)) {
            console.error("No such SSN node: %o", address);
        } else {
            const arguments = implState.ssnlist[address].arguments;
            const rewards = arguments[2];
            console.log("Stake rewards for %o: %o", address, rewards);
        }

    } catch (err) {
        console.log(err);
    }
    console.log("------------------------ end get stake rewards amount ------------------------\n");
}

main();