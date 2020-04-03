# Getting Started
This is a collection of Javascript scripts that shows how to interact with the staking smart contract from deploying, deposit stake deposit to withdraw rewards. 

These scripts utilizes [Zilliqa-JavaScript-Library](https://github.com/Zilliqa/Zilliqa-JavaScript-Library) and has been tested on **NodeJS 12.9.0**.

# Contents
The scripts are divided into folders based on the following roles
- Contract admin
- Verifier
- Staked seed node operator
- Proxy contract admin

# Requirements
- NodeJS >= 10.6

# Installation
Install dependencies via `npm`:
```
cd javascript
npm install
```

# Configuration
Open the `.js` file with a text editor which you want to execute and edit the parameters accordingly:
```
const API = 'http://localhost:5555' // use https://dev-api.zilliqa.com for Dev Testnet
const CHAIN_ID = 1; // use 333 for developer testnet
const PRIVATE_KEY = 'd96e9eb5b782a80ea153c937fa83e5948485fbfc8b7e7c069d7b914dbc350aba';
...
...
```

# Execution
Execute a script to interact with the smart contract:
```
node deploy.js
```

