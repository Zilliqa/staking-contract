# Getting Started
This is a collection of Javascript scripts that shows how to interact with the staking smart contract from deploying, deposit stake deposit to withdraw rewards. These scripts utilizes `Zilliqa-JavaScript-Library` and tested on **NodeJS 12.9.0**.

Install dependencies via `npm`:
```
cd javascript
npm install
```

Open the `.js` file with a text editor which you want to execute and edit the parameters accordingly:
```
const API = 'http://localhost:5555' // use https://dev-api.zilliqa.com for Dev Testnet
const CHAIN_ID = 1; // use 333 for developer testnet
const PRIVATE_KEY = 'd96e9eb5b782a80ea153c937fa83e5948485fbfc8b7e7c069d7b914dbc350aba';
...
...
```

Execute a script to interact with the smart contract:
```
node deploy.js
```

# Requirements
- NodeJS >= 10.6