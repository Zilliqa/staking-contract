# Seed Node Staking Contracts

Seed node staking in [Zilliqa](https://scilla.readthedocs.io/en/latest/) as described in [ZIP-3](https://github.com/Zilliqa/ZIP/blob/master/zips/zip-3.md) makes use of three contracts written in [Scilla](https://scilla.readthedocs.io/en/latest/). This repository is the central portal that collates together the contracts, documentations around them and unit tests to deploy and invoke the contracts on the network.

In the sections below, we describe in detail: 1) the purpose of each contract, 2) their structure and specifications, 3) running unit tests for the contracts.

## Smart Contract Specifications

The table below summarizes the three contracts that ZIP-3 will broadly use:


| Contract Name | File and Location | Description |
|--|--| --|
|SSNList| [`ssnlist.scilla`](./contracts/ssnlist.scilla)  | The main contract that keeps track of Staked Seed Nodes _aka_ SSNs and the amount staked and the rewards.|
|SSNListProxy| [`proxy.scilla`](./contracts/proxy.scilla)  | A proxy contract that sits on top of the SSNList contract. Any call to the SSNList contract must come from SSNListProxy. This contracts facilitates upgradeability of SSNList contract in case a bug is found.|
|Wallet| [`multisig_wallet.scilla`](./contracts/multisig_wallet.scilla)  | A multisig wallet contract tailored to work with the SSNListproxy contract. Certain transitions in the SSNListProxy contract can only be invoked when k-out-of-n users have agreed to do so. This logic is handled using the Wallet contract. |

### SSNList Contract Specification



