# Non-Custodial Seed Node Staking Contracts

Non-custodial seed node staking in [Zilliqa](https://www.zilliqa.com) as
described in
[ZIP-11](https://github.com/Zilliqa/ZIP/blob/zip-11/zips/zip-11.md) makes use
of four contracts written in
[Scilla](https://scilla.readthedocs.io/en/latest/).  This repository is the
central portal that collates together the contracts, documentations around
them, unit tests, and scripts to deploy and run the contracts on the network.

In the sections below, we describe in detail: 1) the purpose of each contract,
2) their structure and specifications, 3) running unit tests for the contracts.

# Table of Content

- [Overview](#overview)
- [SSNList Contract Specification](#ssnlist-contract-specification)
  * [Roles and Privileges](#roles-and-privileges)
  * [Data Types](#data-types)
  * [Immutable Parameters](#immutable-parameters)
  * [Mutable Fields](#mutable-fields)
  * [Transitions](#transitions)
    + [Housekeeping Transitions:](#housekeeping-transitions-)
    + [Pause Transitions](#pause-transitions)
    + [SSN Operation Transitions](#ssn-operation-transitions)
- [SSNListProxy Contract Specification](#ssnlistproxy-contract-specification)
  * [Roles and Privileges](#roles-and-privileges-1)
  * [Immutable Parameters](#immutable-parameters-1)
  * [Mutable Fields](#mutable-fields-1)
  * [Transitions](#transitions-1)
      + [Housekeeping Transitions](#housekeeping-transitions)
      + [Relay Transitions](#relay-transitions)
- [Multi-signature Wallet Contract Specification](#multi-signature-wallet-contract-specification)
  * [General Flow](#general-flow)
  * [Roles and Privileges](#roles-and-privileges-2)
  * [Immutable Parameters](#immutable-parameters-2)
  * [Mutable Fields](#mutable-fields-2)
  * [Transitions](#transitions-2)
    + [Submit Transitions](#submit-transitions)
    + [Action Transitions](#action-transitions)


# Overview

The table below summarizes the purpose of the four contracts that ZIP-11 will
broadly use:

| Contract Name | File and Location | Description |
|--|--| --|
|SSNList| [`ssnlist.scilla`](./contracts/ssnlist.scilla)  | The main contract that keeps track of Staked Seed Nodes _aka_ SSNs, the delegators, the amount staked by a delegator with an SSN, and available rewards, etc.|
|SSNListProxy| [`proxy.scilla`](./contracts/proxy.scilla)  | A proxy contract that sits on top of the SSNList contract. Any call to the `SSNList` contract must come from `SSNListProxy`. This contracts facilitates upgradeability of the `SSNList` contract in case a bug is found.|
|Wallet| [`multisig_wallet.scilla`](./contracts/multisig_wallet.scilla)  | A multisig wallet contract tailored to work with the `SSNListproxy` contract. Certain transitions in the `SSNListProxy` contract can only be invoked when k-out-of-n users have agreed to do so. This logic is handled using the `Wallet` contract. |
|gZILToken| [`gzil.scilla`](./contracts/gzil.scilla)  | A [ZRC-2](https://github.com/Zilliqa/ZRC/blob/master/zrcs/zrc-2.md) compliant fungible token contract. gZIL tokens represent governance tokens. They are issued alongside staking rewards whenever a delegator withdraws her staking rewards. |


# SSNList Contract Specification

The SSNList contract is the main contract that is central to the entire staking
infrastructure. 


## Roles and Privileges

The table below describes the roles and privileges that this contract defines:

| Role | Description & Privileges|                                    
| --------------- | ------------------------------------------------- |
| `ssn`           | A registered SSN that provides the seed node service and gets rewarded for the service. |
| `verifier`      | An entity that checks the health of SSNs and rewards them accordingly for their service.                                 |
| `admin`         | The administrator of the contract.  `admin` is a multisig wallet contract (i.e., an instance of `Wallet`).    |
| `initiator`     | The user who calls the `SSNListProxy` that in turns call the `SSNList` contract. |
| `delegator`     | A token holder who wishes to delegate her tokens to an SSN for staking. The delegator earns a portion of the reward that the SSN receives. |

## Data Types

The contract defines and uses several custom ADTs that we describe below: 


1. SSN Data Type: 

```ocaml
type Ssn =
| Ssn of Bool Uint128 Uint128 String String String Uint128 Uint128 Uint128 ByStr20
```

```ocaml
(* Each SSN has the following fields: *)

(* ActiveStatus      : Bool *)
(*                     Represents whether the SSN has the minimum stake amount and therefore ready to participate in staking and receive rewards. *)
(* StakeAmount       : Uint128 *)
(*                     Total stake that can be used for reward calculation. *)
(* StakeRewards      : Uint128 *)
(*                     (Unwithdrawn) Reward accumulated so far across all cycles. It only includes the reward that the SSN can distribute to its delegators. It does not include SSN's own commission. *)
(* Name              : String *)
(*                     A human-readable name for this SSN. *)
(* URLRaw            : String *)
(*                     Represents "ip:port" of the SSN serving raw API requests. *)
(* URLApi            : String *)
(*                     Representing URL exposed by SSN serving public API requests. *)
(* BufferedDeposit   : Uint128 *)
(*                     Stake deposit that cannot be counted as a part of reward calculation for the ongoing reward cycle. But, to be considered for the next one. *)
(* Commission        : Uint128 *)
(*                     Percentage of incoming rewards that the SSN takes. *)
(* CommissionRewards : Uint128 *)
(*                     Number of ZILs earned as commission by the SSN. *)
(* ReceivingAddress   : ByStr20 *)
(*                     Address to be used to receive commission. *)
```

2. SSNRewardShare Data Type:

```ocaml
type SsnRewardShare =
| SsnRewardShare of ByStr20 Uint128
``` 

```ocaml
(* SSNRewardShare has the following fields: *)

(*  SSNAddress        : ByStr20 *)
(*                      Address of the SSN. *)
(*  RewardShare       : Uint128 *)
(*                      This is the integer representation of the reward assigned by the verifier to this SSN for this cycle. *) 
                        It's floor(NumberOfDSEpochsInCurrentCycle * 110,000 * VerificationPassed) *)
```

3. DelegCycleInfo Data Type:

```ocaml
type DelegCycleInfo =
| DelegCycleInfo of ByStr20 Uint128 ByStr20
```

```ocaml
(*    Each DelegCycleInfo has the following fields: *)
(*    SSNAddress          : ByStr20 *)
(*                          Address of the SSN. *)
(*    StakeDuringTheCycle : Uint128 *)
(*                          Represents the amount staked during this cycle for the given SSN. *)
(*    DelegAddress        : ByStr20 *)
(*                          Address of Delegator. *)
```

4. SSNCycleInfo Data Type:

```ocaml
type SSNCycleInfo =
| SSNCycleInfo of Uint128 Uint128
```
```ocaml
(*   Each SSNCycleInfo has the following fields: *)
(*   TotalStakeDuringTheCycle            : Uint128 *)
(*                                          Represents the amount staked during this cycle for the given SSN. *)
(*    TotalRewardEarnedDuringTheCycle    : Uint128 *)
(*                                         Represents the total reward earned during this cycle for the given SSN. *)
```

5. Error Data Type:

```ocaml
type Error =
  | ContractFrozenFailure (* Contract is paused *)
  | VerifierValidationFailed (* Initiator is not verifier *)
  | AdminValidationFailed (* Initiator is not admin *)
  | ProxyValidationFailed (* Caller is not proxy *)
  | DelegDoesNotExistAtSSN (* Delegator does not exist at the given SSN *)
  | UpdateStakingParamError (* Error in updating staking parameter *)
  | DelegHasBufferedDeposit (* Delegator has some buffered deposit. *)
  | ChangeCommError (* Commission could not be changed *)
  | SSNNotExist (* SSN does not exist *)
  | DelegNotExist (* Delegator does not exist *)
  | SSNAlreadyExist (* SSN already exists *)
  | DelegHasUnwithdrawnRewards (* Delegator has unwithdrawn Rewards *)
  | DelegHasNoSufficientAmt (* Delegator does not have sufficient amount to withdraw *)
  | SSNNoComm (* SSN has no commission left to withdraw *)
  | DelegStakeNotEnough (* Delegator's stake is not above minimum *)
```

## Immutable Parameters

The table below lists the parameters that are defined at the contract deployment time and hence cannot be changed later on.

| Name | Type | Description |                                    
| --------------- | ------------------------------------------------- |-|
| `init_admin` | `ByStr20` | The initial admin of the contract.          |
| `proxy_address` | `ByStr20` | Address of the `SSNListProxy` contract.  |

## Mutable Fields

The contract defines and uses a custom ADT named `Ssn` as explained below: 

```
type Ssn = 
| Ssn of Bool Uint128 Uint128 String String Uint128
(* The first argument of type Bool represents the status of the SSN. An active SSN will have this field set to True.  *)
(* The second argument of type Uint128 represents the amount staked by the SSN. *)
(* The third argument of type Uint128 represents the rewards that this SSN can withdraw. *)
(* The fourth argument of type String represents the raw URL for the SSN to fetch raw blockchain data. *)
(* The fifth argument of type String represents the URL API endpoint for the SSN similar to api.zilliqa.com. *)
(* The sixth and the last argument of type Uint128 represents the deposit made by the SSN that cannot be considered for reward calculations in the current reward cycle. *)
```

The table below presents the mutable fields of the contract and their initial values. 

| Name        | Type       | Initial Value                           | Description                                        |
| ----------- | --------------------|--------------- | -------------------------------------------------- |
| `ssnlist`   | `Map ByStr20 Ssn` | `Emp ByStr20 Ssn` |Mapping between SSN addresses and the corresponding `Ssn` information. |
| `verifier`   | `Option ByStr20` | `None {ByStr20}` | The address of the `verifier`. |
| `minstake`  | `Uint128` | `Uin128 0`       | Minimum stake required to activate an SSN (in `Qa`, where `1 Qa = 10^-12 ZIL`). |
| `maxstake`  | `Uint128`  | `Uint128 0`.                       | Maximum stake (in `Qa`) allowed for each SSN. |
| `contractmaxstake`  | `Uint128`  | `Uint128 0` | The maximum amount (in `Qa`) that can ever be staked across all SSNs. |
| `totalstakeddeposit`  | `Uint128`  | `Uint128 0` | The total amount (in `Qa`) that is currently staked in the contract. |
| `contractadmin` | `ByStr20` |  `init_admin` | Address of the administrator of this contract. |
|`paused` | `ByStr20` | `True` | A flag to record the paused status of the contract. Certain transitions in the contract cannot be invoked when the contract is paused. |
|`lastrewardblocknum` | `Uint32` | `Uint32 0` | The block number when the last reward was distributed. |

## Transitions 

Note that each of the transitions in the `SSNList` contract takes `initiator` as a parameter which as explained above is the caller that calls the `SSNListProxy` contract which in turn calls the `SSNList` contract. 

> Note: No transition in the `SSNList` contract can be invoked directly. Any call to the `SSNList` contract must come from the `SSNListProxy` contract.

All the transitions in the contract can be categorized into three categories:

* **Housekeeping Transitions:** Meant to facilitate basic admin-related tasks.
* **Pause Transitions:** Meant to pause and un-pause the contract.
* **SSN Operation Transitions:** The core transitions that the `verifier` and the SSNs will invoke as a part of the SSN operation.

Each of these category of transitions are presented in further detail below.

### Housekeeping Transitions:

| Name        | Params     | Description | Callable when paused?|
| ----------- | -----------|-------------|:--------------------------:|
| `update_admin` | `admin : ByStr20, initiator : ByStr20` | Replace the current `contractadmin` by `admin`. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: |
| `update_verifier` | `verif : ByStr20, initiator : ByStr20` | Replace the current `verifier` by `verif`. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: |
| `update_minstake` | `min_stake : Uint128, initiator : ByStr20` | Update the value of the field `min_stake` to the input value `min_stake`. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| <center>:x:</center> |
| `update_maxstake` | `max_stake : Uint128, initiator : ByStr20` | Update the value of the field `max_stake` to the input value `max_stake`. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| <center>:x:</center> |
| `update_contractmaxstake` | `max_stake : Uint128, initiator : ByStr20` | Update the value of the field `contractmaxstake` to the input value `max_stake`. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| <center>:x:</center> |
| `drain_contract_balance` | `initiator : ByStr20` | Allows the admin to withdraw the entire balance of the contract. It should only be invoked in case of emergency. The withdrawn ZILs go to a multsig wallet contract that represents the `admin`. :warning: **Note:** `initiator` must be the current `contractadmin` of the contract. | :heavy_check_mark:|


### Pause Transitions

| Name        | Params     | Description | Callable when paused?|
| ----------- | -----------|-------------|:--------------------------:|
| `pause` | `initiator : ByStr20`| Pause the contract temporarily to stop any critical transition from being invoked. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: | 
| `unpause` | `initiator : ByStr20`| Un-pause the contract to re-allow the invocation of all transitions. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: |


### SSN Operation Transitions

| Name        | Params     | Description | Callable when paused?|
| ----------- | -----------|-------------|:--------------------------:|
| `add_ssn` | `ssnaddr : ByStr20, stake_amount : Uint128, rewards : Uint128, urlraw : String, urlapi : String, buffered_deposit : Uint128, initiator : ByStr20`| Add a new SSN with the passed values. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | <center>:x:</center> | 
| `remove_ssn` | `ssnaddr : ByStr20, initiator : ByStr20`| Remove a given SSN with address `ssnaddr`. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | <center>:x:</center> |
| `stake_deposit` | `initiator : ByStr20` | Accept the deposit of ZILs from the `initiator` which should be a registered SSN in the contract. | <center>:x:</center> | 
| `assign_stake_reward` | `ssnreward_list : List SsnRewardShare, reward_blocknum : Uint32, initiator : ByStr20` | To assign rewards to the SSNs based on their performance. Performance checks happen off the chain. <br>  :warning: **Note:** `initiator` must be the current `verifier` of the contract. | <center>:x:</center> |
| `withdraw_stake_rewards` | `initiator : ByStr20` | A registered SSN (`initiator`) can call this transition to withdraw its stake rewards. Stake rewards represent the rewards that an SSN has earned based on its performance. If the SSN has already withdrawn its deposit, then the SSN is removed. | <center>:x:</center> |
| `withdraw_stake_amount` | `amount : Uint128, initiator : ByStr20` | A registered SSN (`initiator`) can call this transition to withdraw its deposit. The amount to be withdrawn is the input value `amount`. Stake amount represents the deposit that an SSN has made so far.   <br> :warning: **Note:** <ul><li> Any partial withdrawal should ensure that the remaining deposit is greater than `min_stake`. Partial withdrawals that push the deposit amount to be lower than `min_stake` are denied. </li> <li> In case there is a non-zero buffereddeposit, withdrawal is denied. </li> <li> If the withdrawal is for the entire deposit and if the rewards have also been withdrawn, then the SSN is removed. </li> </ul>  | <center>:x:</center> |
| `deposit_funds` | `initiator : ByStr20` | Deposit ZILs to the contract.  | <center>:x:</center> |

>**Note:** Ssn custom ADT has a field status of type `Bool`: `True` representing an active SSN. Value of this field depends on whether or not the amount staked by this SSN and the rewards held by this SSN are zero. The table below presents the different possible configurations and the value of the status field for each configuration. Note that when both the staked amount is zero and the reward is zero, the SSN is removed. <br>
>| Staked Amount Value | Reward Amount Value | SSN Status |
>| -- | -- | -- |
>| Zero | Non-zero | `False`|
>| Non-zero | Zero | `True`|
>| Non-zero | Non-zero | `True`|
>| Zero | Zero | `Removed`|

# SSNListProxy Contract Specification

`SSNListProxy` contract is a relay contract that redirects calls to it to the `SSNList` contract.

## Roles and Privileges

The table below describes the roles and privileges that this contract defines:

| Role | Description & Privileges|                                    
| --------------- | ------------------------------------------------- |
| `init_admin`           | The initial admin of the contract which is usually the creator of the contract. `init_admin` is also the initial value of admin. |                                 |
| `admin`    | Current `admin` of the contract initialized to `init_admin`. Certain critical actions can only be performed by the `admin`, e.g., changing the current implementation of the `SSNList` contract. |
|`initiator` | The user who calls the `SSNListProxy` contract that in turns call the `SSNList` contract. |

## Immutable Parameters

The table below lists the parameters that are defined at the contrat deployment time and hence cannot be changed later on.

| Name | Type | Description |
|--|--|--|
|`init_implementation`| `ByStr20` | The address of the `SSNList` contract. |
|`init_admin`| `ByStr20` | The address of the admin. |

## Mutable Fields

The table below presents the mutable fields of the contract and their initial values.

| Name | Type | Initial Value |Description |
|--|--|--|--|
|`implementation`| `ByStr20` | `init_implementation` | Address of the current implementation of the `SSNList` contract. |
|`admin`| `ByStr20` | `init_owner` | Current `admin` of the contract. |


## Transitions

All the transitions in the contract can be categorized into two categories:
- **Housekeeping Transitions** meant to facilitate basic admin related tasks.
- **Relay Transitions** to redirect calls to the `SSNList` contract.

### Housekeeping Transitions

| Name | Params | Description |
|--|--|--|
|`upgradeTo`| `newImplementation : ByStr20` |  Change the current implementation address of the `SSNList` contract. <br> :warning: **Note:** Only the `admin` can invoke this transition|
|`changeProxyAdmin`| `newAdmin : ByStr20` |  Change the current `admin` of the contract. <br> :warning: **Note:** Only the `admin` can invoke this transition.|

### Relay Transitions


These transitions are meant to redirect calls to the corresponding `SSNList` contract and hence their names have an added prefix `proxy`. While, redirecting the contract prepares the `initiator` value that is the address of the caller of the `SSNListProxy` contract. The signature of transitions in the two contracts is exactly the same expect the added last parameter `initiator` for the `SSNList` contract.

| Transition signature in the `SSNListProxy` contract  | Target transition in the `SSNList` contract |
|--|--|
|`pause()` | `pause(initiator : ByStr20)` |
|`unpause()` | `unpause(initiator : ByStr20)` |
|`update_admin(admin: ByStr20)` | `update_admin(admin: ByStr20, initiator : ByStr20)`|
|`update_verifier(verif : ByStr20)` | `update_verifier (verif : ByStr20, initiator: ByStr20)`|
|`drain_contract_balance()` | `drain_contract_balance(initiator : ByStr20)`|
|`update_minstake (min_stake : Uint128)` | `update_minstake (min_stake : Uint128, initiator : ByStr20)`|
|`update_maxstake (min_stake : Uint128)` | `update_maxstake (max_stake : Uint128, initiator : ByStr20)`|
|`update_contractmaxstake (max_stake : Uint128)` | `update_contractmaxstake (max_stake : Uint128, initiator : ByStr20)`|
|`add_ssn (ssnaddr : ByStr20, stake_amount : Uint128, rewards : Uint128, urlraw : String, urlapi : String, buffered_deposit : Uint128)` | `add_ssn (ssnaddr : ByStr20, stake_amount : Uint128, rewards : Uint128, urlraw : String, urlapi : String, buffered_deposit : Uint128, initiator : ByStr20)`|
|`remove_ssn (ssnaddr : ByStr20)` | `remove_ssn (ssnaddr : ByStr20, initiator: ByStr20)`|
|`stake_deposit()` | `stake_deposit (initiator: ByStr20)`|
|`assign_stake_reward (ssnreward_list : List SsnRewardShare, reward_blocknum : Uint32)` | `assign_stake_reward (ssnreward_list : List SsnRewardShare, reward_blocknum : Uint32, initiator: ByStr20)`|
|`withdraw_stake_rewards()` | `withdraw_stake_rewards (initiator : ByStr20)`|
|`withdraw_stake_amount (amount : Uint128)` | `withdraw_stake_amount (amount : Uint128, initiator: ByStr20)`|
|`deposit_funds()` | `deposit_funds (initiator : ByStr20)`|

# Multi-signature Wallet Contract Specification

This contract has two main roles. First, it holds funds that can be paid out to arbitrary users, provided that enough people from a pre-defined set of owners have signed off on the payout.

Second, and more generally, it also represents a group of users that can invoke a transition in another contract only if enough people in that group have signed off on it. In the staking context, it represents the `admin` in the `SSNList` contract. This provides added security for the privileged `admin` role.

## General Flow

Any transaction request (whether transfer of payments or invocation of a foreign transition) must be added to the contract before signatures can be collected. Once enough signatures are collected, the recipient (in case of payments) and/or any of the owners (in the general case) can ask for the transaction to be executed.

If an owner changes his mind about a transaction, the signature can be revoked until the transaction is executed.

This wallet does not allow adding or removing owners, or changing the number of required signatures. To do any of those, perform the following steps:

1. Deploy a new wallet with `owners` and `required_signatures` set to the new values. `MAKE SURE THAT THE NEW WALLET HAS BEEN SUCCESFULLY DEPLOYED WITH THE CORRECT PARAMETERS BEFORE CONTINUING!`
2. Invoke the `SubmitTransaction` transition on the old wallet with the following parameters:
   - `recipient` : The `address` of the new wallet
   - `amount` : The `_balance` of the old wallet
   - `tag` : `AddFunds`
3. Have (a sufficient number of) the owners of the old contract invoke the `SignTransaction` transition on the old wallet. The parameter `transactionId` should be set to the `Id` of the transaction created in step 2.
4. Have one of the owners of the old contract invoke the `ExecuteTransaction` transition on the old contract. This will cause the entire balance of the old contract to be transferred to the new wallet. Note that no un-executed transactions will be transferred to the new wallet along with the funds.

> WARNING: If a sufficient number of owners lose their private keys, or for any other reason are unable or unwilling to sign for new transactions, the funds in the wallet will be locked forever. It is therefore a good idea to set required_signatures to a value strictly less than the number of owners, so that the remaining owners can retrieve the funds should such a scenario occur.
<br> <br> If an owner loses his private key, the remaining owners should move the funds to a new wallet (using the workflow described above) to  ensure that funds are not locked if another owner loses his private key. The owner who originally lost his private key can generate a new key, and the corresponding address be added to the new wallet, so that the same set of people own the new wallet.

## Roles and Privileges

The table below list the different roles defined in the contract.

| Name | Description & Privileges |
|--|--|
|`owners` | The users who own this contract. |

## Immutable Parameters

The table below lists the parameters that are defined at the contract deployment time and hence cannot be changed later on.

| Name | Type | Description |
|--|--|--|
|`owners_list`| `List ByStr20` | List of initial owners. |
|`required_signatures`| `Uint32` | Minimum amount of signatures to execute a transaction. |

## Mutable Fields

The table below presents the mutable fields of the contract and their initial values.

| Name | Type | Initial Value | Description |
|--|--|--|--|
|`owners`| `Map ByStr20 Bool` | `owners_list` | Map of owners. |
|`transactionCount`| `Uint32` | `0` | The number of of transactions  requests submitted so far. |
|`signatures`| `Map Uint32 (Map ByStr20 Bool)` | `Emp Uint32 (Map ByStr20 Bool)` | Collected signatures for transactions by transaction ID. |
|`signature_counts`| `Map Uint32 Uint32` | `Emp Uint32 Uint32` | Running count of collected signatures for transactions. |
|`transactions`| `Map Uint32 Transaction` | `Emp Uint32 Transaction` | Transactions that have been submitted but not exected yet. |

## Transitions

All the transitions in the contract can be categorized into three categories:
- **Submit Transitions:** Create transactions for future signoff.
- **Action Transitions:** Let owners sign, revoke or execute submitted transactions.
- The `_balance` field keeps the amount of funds held by the contract and can be freely read within the implementation. `AddFunds transition` are used for adding native funds(ZIL) to the wallet from incoming messages by using `accept` keyword.

### Submit Transitions

The first transition is meant to submit request for transfer of native ZILs while the other are meant to submit a request to invoke transitions in the `SSNListProxy` contract.

| Name | Params | Description |
|--|--|--|
|`SubmitNativeTransaction`| `recipient : ByStr20, amount : Uint128, tag : String` | Submit a request for transafer of native tokens for future signoffs. |
|`SubmitCustomUpgradeToTransaction`| `proxyContract : ByStr20, newImplementation : ByStr20` | Submit a request to invoke the `upgradeTo` transition in the `SSNListProxy` contract. |
|`SubmitCustomChangeProxyAdminTransaction`| `proxyContract : ByStr20, newAdmin : ByStr20` | Submit a request to invoke the `changeProxyAdmin` transition in the `SSNListProxy` contract. |
|`SubmitCustomUpdateAdminTransaction`| `proxyContract : ByStr20, admin : ByStr20` | Submit a request to invoke the `update_admin` transition in the `SSNListProxy` contract. |
|`SubmitCustomUpdateVerifierTransaction`| `proxyContract : ByStr20, verif : ByStr20` | Submit a request to invoke the `update_verifier` transition in the `SSNListProxy` contract. |
|`SubmitCustomUpdateMinStakeTransaction`| `proxyContract : ByStr20, min_stake : Uint128` | Submit a request to invoke the `update_minstake` transition in the `SSNListProxy` contract. |
|`SubmitCustomUpdateMaxStakeTransaction`| `proxyContract : ByStr20, max_stake : Uint128` | Submit a request to invoke the `update_maxstake` transition in the `SSNListProxy` contract. |
|`SubmitCustomUpdateContractMaxStakeTransaction`| `proxyContract : ByStr20, max_stake : Uint128` | Submit a request to invoke the `update_contractmaxstake` transition in the `SSNListProxy` contract. |
|`SubmitCustomDrainContractBalanceTransaction`| `proxyContract : ByStr20` | Submit a request to invoke the `drain_contract_balance` transition in the `SSNListProxy` contract. |
|`SubmitCustomAddSsnTransaction`| `proxyContract : ByStr20, ssnaddr : ByStr20, stake_amount : Uint128, rewards : Uint128, urlraw : String, urlapi : String, buffered_deposit : Uint128` | Submit a request to invoke the `add_ssn` transition in the `SSNListProxy` contract. |
|`SubmitCustomRemoveSsnTransaction`| `proxyContract : ByStr20, ssnaddr : ByStr20` | Submit a request to invoke the `remove_ssn` transition in the `SSNListProxy` contract. |



### Action Transitions

| Name | Params | Description |
|--|--|--|
|`SignTransaction`| `transactionId : Uint32` | Sign off on an existing transaction. |
|`RevokeSignature`| `transactionId : Uint32` | Revoke signature of an existing transaction, if it has not yet been executed. |
|`ExecuteTransaction`| `transactionId : Uint32` | Execute signed-off transaction. |
