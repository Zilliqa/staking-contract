## Testing using Zilliqa Isolated Server
The Zilliqa isolated server is a simulated environment for smart contract testing. It uses the `accountstore` and scilla interpreter to run smart contract transactions in the absense of consensus protocol. This enables more rapid testing for smart contracts.

### Requirements

Golang (minimum version: go.1.12):
* [Download page](https://golang.org/dl/)
* [Installation instructions](https://golang.org/doc/install)

Zilliqa command tool [zli](https://github.com/Zilliqa/zli)

### Private keys for unit tests

We are using an isolated server to run our tests, some of the tests require a large amount of ZILs to be executed. Instead of letting you generate your own private keys, we have prepared two private keys for you and will send them privately to you.

Let's name these two private keys as `pri1` and `pri2` in the following context.

### Setup zli command tool

#### 1. Init wallet configuration with specific private key

Run the following command to init your wallet with the provided private key:

```shell script
zli wallet from -p `pri1`
```

A `.zilliqa` file containing the wallet configuration will be generated under your `USERS` directory after running this. You can checkout this either by using `cat ~/.zilliqa` or `zli wallet echo`, you will see something like the following:

```json
{"api":"https://dev-api.zilliqa.com/","chain_id":333,"default_account":{"private_key":"e53d1c3edaffc7a7bab5418eb836cf75819a82872b4a1a0f1c7fcf5c3e020ccc","public_key":"036695e20c8339bd3aab70aead5fc0e35ade557b4d00f0552c62afa220ad0ee149","address":"ad7d96b8b4d7a13b96b0dd1081832606090c096d","bech_32_address":"zil1447edw9567snh94sm5ggrqexqcyscztddt2t94"},"accounts":[{"private_key":"e53d1c3edaffc7a7bab5418eb836cf75819a82872b4a1a0f1c7fcf5c3e020ccc","public_key":"036695e20c8339bd3aab70aead5fc0e35ade557b4d00f0552c62afa220ad0ee149","address":"ad7d96b8b4d7a13b96b0dd1081832606090c096d","bech_32_address":"zil1447edw9567snh94sm5ggrqexqcyscztddt2t94"}]}
```


#### 2. Modify `api` and `chain_id`

Since we are using the `isolated server` rather than `community devnet`, we need to rewrite the parameters `api` and `chain_id` in `.zilliqa`:

    "api": "https://zilliqa-isolated-server.zilliqa.com/",
    "chain_id": 1,

The changed file should looked similar to the following:

```json
{"api":"https://zilliqa-isolated-server.zilliqa.com/","chain_id":1,"default_account":{"private_key":"e53d1c3edaffc7a7bab5418eb836cf75819a82872b4a1a0f1c7fcf5c3e020ccc","public_key":"036695e20c8339bd3aab70aead5fc0e35ade557b4d00f0552c62afa220ad0ee149","address":"ad7d96b8b4d7a13b96b0dd1081832606090c096d","bech_32_address":"zil1447edw9567snh94sm5ggrqexqcyscztddt2t94"},"accounts":[{"private_key":"e53d1c3edaffc7a7bab5418eb836cf75819a82872b4a1a0f1c7fcf5c3e020ccc","public_key":"036695e20c8339bd3aab70aead5fc0e35ade557b4d00f0552c62afa220ad0ee149","address":"ad7d96b8b4d7a13b96b0dd1081832606090c096d","bech_32_address":"zil1447edw9567snh94sm5ggrqexqcyscztddt2t94"}]}
```

### Run tests

With `zli` installed and configured, we can use `go run <pri1> <pri2> <api>` to run tests, for instance:

```go
go run main/main.go e53d1c3edaffc7a7bab5418eb836cf75819a82872b4a1a0f1c7fcf5c3e020ccc 21d1af225dab7b791d656f2fba2f8b5ca513327f9dce29473a7ec2aba5351318 https://zilliqa-isolated-server.zilliqa.com/
```
