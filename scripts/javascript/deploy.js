/*
 * deploy
 */
const fs = require('fs');
const path = require('path');
const { BN, Long, bytes, units } = require('@zilliqa-js/util');
const { Zilliqa } = require('@zilliqa-js/zilliqa');
const { toBech32Address, getAddressFromPrivateKey } = require('@zilliqa-js/crypto');

const zilliqa = new Zilliqa('https://dev-api.zilliqa.com');
const CHAIN_ID = 2;
const MSG_VERSION = 1;
const VERSION = bytes.pack(CHAIN_ID, MSG_VERSION);

const PRIVATE_KEY = '';
const GAS_PRICE = units.toQa('1000', units.Units.Li);

const PROXY_CONTRACT_PATH = "./proxy.scilla";
const SSNLIST_CONTRACT_PATH = "./ssnlist.scilla";

async function main() {
    zilliqa.wallet.addByPrivateKey(PRIVATE_KEY);
    const address = getAddressFromPrivateKey(PRIVATE_KEY);

    console.log("Invoking deploy...");
    console.log("Your account address is:");
    console.log(`${address}`);

    try {
        // console.log(`Deploying proxy contract....` + PROXY_CONTRACT_PATH);
        // const code = fs.readFileSync(path.join(__dirname, PROXY_CONTRACT_PATH), 'ascii');
    
        // const init = [
        //   // this parameter is mandatory for all init arrays
        //   {
        //       vname: "_scilla_version",
        //       type: "Uint32",
        //       value: "0",
        //   },
        //   {
        //       vname: "init_admin",
        //       type: "ByStr20",
        //       value: `${address}`
        //   },
        //   {
        //       vname: "init_implementation",
        //       type: "ByStr20",
        //       value: `${address}`
        //   }
        // ];
        // const contract = zilliqa.contracts.new(code, init);
        // const [deployTx, proxy] = await contract.deploy(
        //     {
        //         version: VERSION,
        //         gasPrice: GAS_PRICE,
        //         gasLimit: Long.fromNumber(10000)
        //     },
        //     33,
        //     1000,
        //     false
        // );

        // check the pending status
        // const pendingStatus = await zilliqa.blockchain.getPendingTxn(deployTx.id);
        // console.log(`Pending status is: `);
        // console.log(pendingStatus.result);

        // process confirm
        // console.log(`The transaction id is:`, deployTx.id);
        // console.log('The contract address is: %o', proxy.address);

        // deploy ssnlist contract
        console.log(`Deploying a ssnlist contract....` + SSNLIST_CONTRACT_PATH);
        const code2 = fs.readFileSync(path.join(__dirname, SSNLIST_CONTRACT_PATH), 'ascii');
        const init2 = [
            {
                vname: "_scilla_version",
                type: "Uint32",
                value: "0"
            },
            {
                vname: "init_admin",
                type: "ByStr20",
                value: "0x1234567890123456789012345678901234567890"
            },
            {
                vname: "proxy_address",
                type: "ByStr20",
                value: "0x1234567890123456789012345678901234567890"
            }
        ];

        const contract2 = zilliqa.contracts.new(code2, init2);
        const [deployTx2, ssnlist] = await contract2.deploy(
            {
                version: VERSION,
                gasPrice: GAS_PRICE,
                gasLimit: Long.fromNumber(10000)
            },
            33,
            1000,
            false
        );

        // check the pending status
        // const pendingStatus2 = await zilliqa.blockchain.getPendingTxn(deployTx2.id);
        // console.log(`Pending status is: `);
        // console.log(pendingStatus2.result);

        // process confirm
        console.log(`The transaction id is:`, deployTx2.id);
        console.log('The ssnlist contract address is: %o', ssnlist.address);

        // upgrade proxy contract with ssn implementation

    } catch (err) {
        console.log(err);
    }
}

main();