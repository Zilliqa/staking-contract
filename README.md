# Seed Node Staking Contracts

Seed node staking in [Zilliqa](https://scilla.readthedocs.io/en/latest/) as described in [ZIP-3](https://github.com/Zilliqa/ZIP/blob/master/zips/zip-3.md) makes use of three contracts written in [Scilla](https://scilla.readthedocs.io/en/latest/). This repository is the central portal that collates together the contracts, documentations around them and unit tests to deploy and invoke the contracts on the network.

In the sections below, we describe in detail: 1) the purpose of each contract, 2) their structure and specifications, 3) running unit tests for the contracts.

# Overview

The table below summarizes the three contracts that ZIP-3 will broadly use:

| Contract Name | File and Location | Description |
|--|--| --|
|SSNList| [`ssnlist.scilla`](./contracts/ssnlist.scilla)  | The main contract that keeps track of Staked Seed Nodes _aka_ SSNs and the amount staked and the rewards.|
|SSNListProxy| [`proxy.scilla`](./contracts/proxy.scilla)  | A proxy contract that sits on top of the SSNList contract. Any call to the `SSNList` contract must come from `SSNListProxy`. This contracts facilitates upgradeability of `SSNList` contract in case a bug is found.|
|Wallet| [`multisig_wallet.scilla`](./contracts/multisig_wallet.scilla)  | A multisig wallet contract tailored to work with the `SSNListproxy` contract. Certain transitions in the `SSNListProxy` contract can only be invoked when k-out-of-n users have agreed to do so. This logic is handled using the `Wallet` contract. |

# SSNList Contract Specification

The SSNList contract is the main contract that is central to the entire staking infrastructure. 


## Role and Privileges

The table below describes the roles and privileges that this contract defines:

| Role | Description & Privileges|                                    
| --------------- | ------------------------------------------------- |
| `ssn`           | A registered SSN that provides the seed node service and gets rewarded for the service. |
| `verifier`      | A entity that checks the health of an SSN and rewards them accordingly for their service.                                 |
| `admin`    | The administrator of the contract.  Admin is a multisig wallet contract (aka `Wallet`).    |
|`initiator` | The user who calls the `SSNListProxy` that in turns call the `SSNList` contract. |

## Immutable Parameters

The table below lists the parameters that are defined at the contract deployment time and hence cannot be changed later on.

| Name | Type | Description |                                    
| --------------- | ------------------------------------------------- |-|
| `init_admin` | `ByStr20` | The initial admin of the contract.          |
| `proxy_address` | `ByStr20` | Address of the `SSNListProxy` contract.  |

## Mutable Fields

The contract defines and uses a custom ADT named `Ssn` defined as follows: 

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

The table below presents the mutable fields of the contract and their initial values. 

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

## Transitions 

Note that each of the transitions in the `SSNList` contract takes `initiator` as a parameter which as explained above is the caller that calls the `SSNListProxy` contract which in turn calls the `SSNList` contract. 

> Note: No transition in the `SSNList` contract can be invoked directly. Any call to the `SSNList` contract must come from the `SSNListProxy` contract.

All the transitions in the contract can be categorized into three categories:

* **Housekeeping Transitions:** Meant to facilitate basic admin related tasks.
* **Pause Transitions:** Meant to pause and un-pause the contract.
* **SSN Operation Transitions:** The core transitions that the `verifier` and the SSNs will invoke as a part of the SSN operation.

Each of these category of transitions are presented in further details below.

### Housekeeping Transitions:

| Name        | Params     | Description | Callable when paused?|
| ----------- | -----------|-------------|----------------------------|
| `update_admin` | `admin : ByStr20, initiator : ByStr20` | Replace the current `contractadmin` by `admin`. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| <center>:x:</center> |
| `update_verifier` | `verif : ByStr20, initiator : ByStr20` | Replace the current `verifier` by `verif`. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| <center>:x:</center> |
| `update_minstake` | `min_stake : Uint128, initiator : ByStr20` | Update the value of the field `min_stake` to the input value `min_stake`. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| <center>:x:</center> |
| `update_maxstake` | `max_stake : Uint128, initiator : ByStr20` | Update the value of the field `max_stake` to the input value `max_stake`. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| <center>:x:</center> |
| `update_contractmaxstake` | `max_stake : Uint128, initiator : ByStr20` | Update the value of the field `contractmaxstake` to the input value `max_stake`. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| <center>:x:</center> |
| `drain_contract_balance` | `initiator : ByStr20` | Allows the admin to withdraw the entire balance of the contract. It should only be invoked in case of emergency. The withdrawn ZILs go to a multsig wallet contract that represents the `admin`. :warning: **Note:** `initiator` must be the current `contractadmin` of the contract. | :heavy_check_mark:|


### Pause Transitions

| Name        | Params     | Description | Callable when paused?|
| ----------- | -----------|-------------|----------------------------|
| `pause` | `initiator : ByStr20`| Pause the contract temporarily to stop any critical transition from being invoked <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: | 
| `unpause` | `initiator : ByStr20`| Un-pause the contract to re-allow the invocation of all transitions. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: |


### SSN Operation Transitions

| Name        | Params     | Description | Callable when paused?|
| ----------- | -----------|-------------|----------------------------|
| `add_ssn` | `ssnaddr : ByStr20, stake_amount : Uint128, rewards : Uint128, urlraw : String, urlapi : String, buffered_deposit : Uint128, initiator : ByStr20`| Add a new SSN with the passed values. <br>  :warning: **Note:** `initiator` must be the current `admin` of the contract.  | <center>:x:</center> | 
| `remove_ssn` | `ssnaddr : ByStr20, initiator : ByStr20`| Remove a given SSN with address `ssnaddr`. <br>  :warning: **Note:** `initiator` must be the current `admin` of the contract.  | <center>:x:</center> |
| `stake_deposit` | `initiator : ByStr20` | Accept the deposit of ZILs from the `initiator` which should be a registered SSN in the contract. | <center>:x:</center> | 
| `assign_stake_reward` | `ssnreward_list : List SsnRewardShare, reward_blocknum : Uint32, initiator : ByStr20` | To assign rewards to the SSNs based on their performance. Performance checks happen off-the-chain. <br>  :warning: **Note:** `initiator` must be the current `verifier` of the contract. | <center>:x:</center> |
| `withdraw_stake_rewards` | `initiator : ByStr20` | A registered SSN (`initiator`) can all this transition to withdraw its stake rewards. Stake rewards represent the rewards that an SSN has earned based on its performance. If the SSN has already withdrawn its deposit, then the SSN is removed. | <center>:x:</center> |
| `withdraw_stake_amount` | `amount : Uint128, initiator : ByStr20` | A registered SSN (`initiator`) can all this transition to withdraw its deposit. The amount to be withdrawn is the input value `amount`. Stake amount represents the deposit that an SSN has made so far.   <br> :warning: **Note:** 1) Any partial withdrawal should ensure that the remaining deposit is greater than `min_stake`. Partial withdrawals that push the deposit amount to be lower than `min_stake` are denied, 2) In case there is a non-zero buffereddeposit, withdrawal is denied. 3) If the withdrawal is for the entire deposit and if the rewards have also been withdrawn then remove the SSN.   | <center>:x:</center> |
| `depsit_funds` | `initiator : ByStr20` | Deposit ZILs to the contract.  | <center>:x:</center> |


**Note:** Each SSN type has a field status of type `Bool`. The value of this field is `True` (representing an active SSN) depends on whether or not the amount staked by this SSN and the rewards held by this SSN is zero. The table below presents the different possible configurations and the value of the status field for each configuration. Note that when both the staked amount is zero and the reward is zero, the SSN is removed.

| Staked Amount Value | Reward Amount Value | SSN Status |
| -- | -- | -- |
| Zero | Non-zero | `False`|
| Non-zero | Zero | `True`|
| Non-zero | Non-zero | `True`|
| Zero | Zero | `Removed`|
