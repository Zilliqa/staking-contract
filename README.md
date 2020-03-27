# Seed Node Staking Contracts

Seed node staking in [Zilliqa](https://scilla.readthedocs.io/en/latest/) as described in [ZIP-3](https://github.com/Zilliqa/ZIP/blob/master/zips/zip-3.md) makes use of three contracts written in [Scilla](https://scilla.readthedocs.io/en/latest/). This repository is the central portal that collates together the contracts, documentations around them and unit tests to deploy and invoke the contracts on the network.

In the sections below, we describe in detail: 1) the purpose of each contract, 2) their structure and specifications, 3) running unit tests for the contracts.

# Overview

The table below summarizes the three contracts that ZIP-3 will broadly use:

| Contract Name | File and Location | Description |
|--|--| --|
|SSNList| [`ssnlist.scilla`](./contracts/ssnlist.scilla)  | The main contract that keeps track of Staked Seed Nodes _aka_ SSNs and the amount staked and the rewards.|
|SSNListProxy| [`proxy.scilla`](./contracts/proxy.scilla)  | A proxy contract that sits on top of the SSNList contract. Any call to the SSNList contract must come from SSNListProxy. This contracts facilitates upgradeability of SSNList contract in case a bug is found.|
|Wallet| [`multisig_wallet.scilla`](./contracts/multisig_wallet.scilla)  | A multisig wallet contract tailored to work with the SSNListproxy contract. Certain transitions in the SSNListProxy contract can only be invoked when k-out-of-n users have agreed to do so. This logic is handled using the Wallet contract. |

# SSNList Contract Specification

The SSNList contract is the main contract that is central to the entire staking infrastructure. 


## Role and Privileges

The table below describes the roles and privileges that this contract defines:

| Role | Description & Privileges|                                    
| --------------- | ------------------------------------------------- |
| `ssn`           | A registered SSN that provides the seed node service and gets rewarded for the service. |
| `verifier`      | A entity that checks the health of an SSN and rewards them accordingly for their service.                                 |
| `admin`    | The administrator of the contract.      |


## Immutable Parameters

The table below lists the parameters that are defined at the contract deployment time and hence cannot be changed later on.

| Name | Type | Description |                                    
| --------------- | ------------------------------------------------- |-|
| `init_admin` | `ByStr20` | The iniital admin of the contract.          |
| `proxy_address` | `ByStr20` | Address of the `SSNListProxy` contract.  |

## Mutable Fields

The table below presents the mutable fields of the contract and their initial values. In the table, we use a custom ADT type `Ssn` defined as follows: 

```
type Ssn = 
| Ssn of Bool Uint128 Uint128 String String Uint128
(* The first argument of type Bool represents the status of the SSN. An active SSN will have this field set to be True  *)
(* The second argument of type Uint128 represents the amount staked by the SSN. *)
(* The third argument of type Uint128 represents the rewards that this SSN can withdraw. *)
(* The fourth argument of type String represents the raw URL for the SSN *)
(* The fifth argument of type String represents the URL API endpoint for the SSN *)
(* The sixth and the last argument of type String represents the deposit made by the SSN that cannot be considered for reward calculations. *)
```

| Name        | Type       | Initial Value                           | Description                                        |
| ----------- | --------------------|--------------- | -------------------------------------------------- |
| `ssnlist`   | `Map ByStr20 Ssn` | `Emp ByStr20 Ssn` |Mapping between SSN address and the corresponding `Ssn` information. |
| `verifier`   | `Option ByStr20` | `None {ByStr20}` | The address of the verifier. |
| `minstake`  | `Uint128` | `Uin128 0`       | Minimum `stake_amount` required to activate an SSN. |
| `maxstake`  | `Uint128`  | `Uint128 0`                         | Maximum stake allowed for each SSN. |
| `contractmaxstake`  | `Uint128`  | `Uint128 0` | The maximum amount that can ever be staked across all SSNs. |
| `totalstakeddeposit`  | `Uint128`  | `Uint128 0` | The total amount that currently staked in the contract. |
| `contractadmin` | `ByStr20` |  `init_admin` | Address of the administrator of this contract. |
|`paused` | `ByStr20` | `True` | A flag to record the paused status of the contract. Certain transitions in the contract cannot be invoked when the contract is paused. |
