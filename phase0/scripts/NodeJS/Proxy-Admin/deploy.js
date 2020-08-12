/*
 * deploy
 * deploys the proxy, staking contract, and upgrades the implementation
 */
const fs = require('fs');
const path = require('path');
const { BN, Long, bytes, units } = require('@zilliqa-js/util');
const { Zilliqa } = require('@zilliqa-js/zilliqa');
const { toBech32Address, getAddressFromPrivateKey } = require('@zilliqa-js/crypto');

// change the following parameters
const API = 'http://localhost:5555'
const CHAIN_ID = 1;
const PRIVATE_KEY = 'e53d1c3edaffc7a7bab5418eb836cf75819a82872b4a1a0f1c7fcf5c3e020b89';

const zilliqa = new Zilliqa(API);
const MSG_VERSION = 1;
const VERSION = bytes.pack(CHAIN_ID, MSG_VERSION);
const GAS_PRICE = units.toQa('1000', units.Units.Li);

const PROXY_CONTRACT_PATH = "../../../contracts/proxy.scilla";
const SSNLIST_CONTRACT_PATH = "../../../contracts/ssnlist.scilla";


async function main() {
    try {
        zilliqa.wallet.addByPrivateKey(PRIVATE_KEY);
        const address = getAddressFromPrivateKey(PRIVATE_KEY);
    
        console.log("Your account address is:");
        console.log(`${address}`);
        console.log("------------------------ begin deploy proxy ------------------------\n");
        console.log(`Deploying proxy contract: ` + PROXY_CONTRACT_PATH);
        const code = fs.readFileSync(path.join(__dirname, PROXY_CONTRACT_PATH), 'ascii');
    
        const init = [
          {
              vname: "_scilla_version",
              type: "Uint32",
              value: "0",
          },
          {
              vname: "init_admin",
              type: "ByStr20",
              value: `${address}`
          },
          {
              vname: "init_implementation",
              type: "ByStr20",
              value: `${address}`
          }
        ];
        const contract = zilliqa.contracts.new(code, init);
        const [deployTx, proxy] = await contract.deploy(
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

        // Introspect the state of the underlying transaction
        console.log(`Deployment Transaction ID: ${deployTx.id}`);
        console.log(`Deployment Transaction Receipt`);
        console.log(deployTx.txParams.receipt);
        console.log('proxy address: %o', proxy.address);
        console.log("------------------------ end deploy proxy ------------------------\n");

        // deploy ssnlist contract
        console.log("------------------------ begin deploy ssnlist ------------------------\n");
        console.log(`Deploying a ssnlist contract: %o`, SSNLIST_CONTRACT_PATH);
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
                value: `${address}`
            },
            {
                vname: "proxy_address",
                type: "ByStr20",
                value: `${proxy.address}`
            }
        ];

        const contract2 = zilliqa.contracts.new(code2, init2);
        const [deployTx2, ssnlist] = await contract2.deploy(
            {
                version: VERSION,
                amount: new BN(0),
                gasPrice: GAS_PRICE,
                gasLimit: Long.fromNumber(30000)
            },
            33,
            1000,
            true
        );

        console.log(`Deployment ssnlist transaction ID: ${deployTx2.id}`);
        console.log(`Deployment ssnlist transaction receipt:`);
        console.log(deployTx2.txParams.receipt);
        console.log('ssnlist contract address: %o', ssnlist.address);
        console.log("------------------------ end deploy ssnlist ------------------------\n");

        // upgrade proxy contract with ssn implementation
        console.log("------------------------ start upgrade ------------------------\n");
        const callTx = await proxy.call(
            'upgradeTo',
            [
                {
                    vname: 'newImplementation',
                    type: 'ByStr20',
                    value: `${ssnlist.address}`
                }
            ],
            {
                version: VERSION,
                amount: new BN(0),
                gasPrice: GAS_PRICE,
                gasLimit: Long.fromNumber(30000)
            },
            33,
            1000,
            true
        );
        console.log("transaction: %o", callTx.id);
        console.log(JSON.stringify(callTx.receipt, null, 4));
        console.log("------------------------ end upgrade ------------------------\n");
    } catch (err) {
        console.log(err);
    }
}

main();