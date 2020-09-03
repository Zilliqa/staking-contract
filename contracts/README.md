# Non-Custodial Seed Node Staking Contracts

Non-custodial seed node staking in [Zilliqa](https://www.zilliqa.com) as
described in
[ZIP-11](https://github.com/Zilliqa/ZIP/blob/master/zips/zip-11.md) makes use
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
(*                     Percentage of incoming rewards that the SSN takes. Represented as an integer. If the commission is 10.5%, then it is multiplied by 10^7 and then resulting integer is set as commission. The assumption is that the percentage is up to 7 decimal places. *)
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
(*                                         Represents the amount staked during this cycle for the given SSN. *)
(*   TotalRewardEarnedDuringTheCycle     : Uint128 *)
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
  | DelegHasBufferedDeposit (* Delegator has some buffered deposit. *)
  | ChangeCommError (* Commission could not be changed *)
  | SSNNotExist (* SSN does not exist *)
  | SSNAlreadyExist (* SSN already exists *)
  | DelegHasUnwithdrawnRewards (* Delegator has unwithdrawn Rewards *)
  | DelegHasNoSufficientAmt (* Delegator does not have sufficient amount to withdraw *)
  | SSNNoComm (* SSN has no commission left to withdraw *)
  | DelegStakeNotEnough (* Delegator's stake is not above minimum *)
  | ExceedMaxChangeRate (* SSN is trying to modify the commission rate by over 1% *)
  | ExceedMaxCommRate (* SSN is trying to set the commission rate greater than the allowed max *)
  | InvalidTotalAmt (* Error when the total stake amount is being decreased by an illegal amount *)
  | InvalidRecvAddr (* Commission is being withdrawn to an address different from that of the receiving address *)
  | VerifierNotSet (* Verifier's address is not set in the field *)
  | VerifierRecvAddrNotSet (* Verifier's reward address address is not set in the field *)
```

## Immutable Parameters

The table below lists the parameters that are defined at the contract deployment time and hence cannot be changed later on.

| Name | Type | Description |                                    
| ---------------      | ----------|-                                         |
| `init_admin`         | `ByStr20` | The initial admin of the contract.       |
| `init_proxy_address` | `ByStr20` | The initial address of the `SSNListProxy` contract.  |
| `gzil_address`       | `ByStr20` | Address of the `gZILToken` contract.  |

## Mutable Fields


The table below presents the mutable fields of the contract and their initial values. 

| Name        | Type       | Initial Value                           | Description                                        |
| ----------- | --------------------|--------------- | -------------------------------------------------- |
| `ssnlist`   | `Map ByStr20 Ssn` | `Emp ByStr20 Ssn` | Mapping between SSN addresses and the corresponding `Ssn` information. |
| `comm_for_ssn`   | `Map ByStr20 (Map Uint128 Uint128)` | `Emp ByStr20 (Map Uint128 Uint128)` | `Map (SSNAddress -> Map (RewardCycleNum -> Commission))` |
| `deposit_amt_deleg`   | `Map ByStr20 (Map ByStr20 Uint128)` | `Emp ByStr20 (Map ByStr20 Uint128)` | `Map (DelegatorAddress -> Map (SSNAddress -> AmoutDelegated))` |
| `ssn_deleg_amt`   | `Map ByStr20 (Map ByStr20 Uint128)` | `Emp ByStr20 (Map ByStr20 Uint128)` | `Map (SSNAddress -> Map (DelegatorAddress -> AmountDelegated))` |
| `buff_deposit_deleg`   | `Map ByStr20 (Map ByStr20 (Map Uint128 Uint128))` | `Emp ByStr20 (Map ByStr20 (Map Uint128 Uint128))` | `Map (DelegatorAddress -> Map (SSNAddress -> Map (RewardCycleNum -> BufferedStakeAmount)))` |
| `direct_deposit_deleg`   | `Map ByStr20 (Map ByStr20 (Map Uint128 Uint128))` | `Emp ByStr20 (Map ByStr20 (Map Uint128 Uint128))` | `Map (DelegatorAddress -> Map (SSNAddress -> Map (RewardCycleNum -> UnBufferedStakeAmount)))` |
| `last_withdraw_cycle_deleg`   | `Map ByStr20 (Map ByStr20  Uint128))` | `Emp ByStr20 (Map ByStr20 Uint128))` | `Map (DelegatorAddress -> Map (SSNAddress -> RewardCycleWhenLastWithdrawn))`. For a new delegator that has never deposited any stake for this SSN, this field will store the reward cycle number during which the delegator successfully deposited its stake. |
| `last_buf_deposit_cycle_deleg`   | `Map ByStr20 (Map ByStr20  Uint128))` | `Emp ByStr20 (Map ByStr20 Uint128))` | `Map (DelegatorAddress -> Map (SSNAddress -> RewardCycleWhenLastDeposited))` |
| `stake_ssn_per_cycle`   | `Map ByStr20 (Map Uint128  SSNCycleInfo))` | `Emp ByStr20 (Map Uint128 SSNCycleInfo))` | `Map (SSNAddress -> Map (RewardCycleNum -> SSNCycleInfo))`. **Note that this data type contains information that corresponds to the end of the cycle when the verifier has distributed the reward. It could therefore be different from the information stored in `ssnlist` particularly its `stake` field which gets updated in the middle of a cycle.** |
|`withdrawal_pending` | `Map ByStr20 (Map BNum Uint128)` | `Emp ByStr20 (Map BNum Uint128)` | `Map (DelegatorAddress -> (BlockNumberWhenRewardWithdrawalRequested -> Amount ))` |
| `bnum_req`  | `Uint128` | `Uin128 24000`       | Bonding period in terms of number of blocks. Set to be equivalent to 14 days. |
|`reward_cycle_list` | `List Uint128` | `Nil {Uint128}` | List of all reward cycles. |
| `verifier`   | `Option ByStr20` | `None {ByStr20}` | The address of the `verifier`. |
| `verifier_receiving_addr`   | `Option ByStr20` | `None {ByStr20}` | The address to receive verifier's rewards. |
| `minstake`  | `Uint128` | `Uin128 10000000000000000000`       | Minimum stake required to activate an SSN (1 mil ZIL expressed in `Qa`, where `1 ZIL = 10^12 Qa`). |
| `mindelegstake`  | `Uint128` | `Uin128 1000000000000000`       | Minimum stake for a delegator (1000 ZIL expressed in Qa where 1 ZIL = 10^12 Qa). |
| `contractadmin` | `ByStr20` |  `init_admin` | Address of the administrator of this contract. |
| `proxyaddr` | `ByStr20` |  `init_proxy_address` | Address of the proxy contract. |
| `gziladdr` | `ByStr20` |  `gzil_address` | Address of the gzil contract. |
|`lastrewardcyle` | `Uint128` | `Uint128 1` | The block number when the last reward was distributed. |
|`paused` | `ByStr20` | `True` | A flag to record the paused status of the contract. Certain transitions in the contract cannot be invoked when the contract is paused. |
| `maxcommchangerate`| `Uint128` | `Uint128 1` | The maximum rate change that an SSN is allowed to make across cycles. Set to 1%. |
| `maxcommrate`| `Uint128` | `Uint128 1000000000` | The maximum commission rate that an SSN can charge. Set to 100% but represented as an integer multiplied by 10^7. |
| `totalstakeamount`  | `Uint128`  | `Uint128 0` | The total amount (in `Qa`) that is currently staked in the contract. It only corresponds to the stake with active SSNs that is unbuffered and therefore can be taken into account for reward calculation. |

## Transitions 

Note that each of the transitions in the `SSNList` contract takes `initiator` as a parameter which as explained above is the caller that calls the `SSNListProxy` contract which in turn calls the `SSNList` contract. 

> Note: No transition in the `SSNList` contract can be invoked directly. Any call to the `SSNList` contract must come from the `SSNListProxy` contract.

All the transitions in the contract can be categorized into four categories:

* **Housekeeping Transitions:** Meant to facilitate basic admin-related tasks.
* **Delegator Transitions:** The transitions that the delegators will invoke as a part of the SSN operation.
* **SSN Transitions:** The transitions that the SSNs will invoke as a part of the SSN operation.
* **Verifier Transitions:** The transitions that the verifier will invoke as a part of the SSN operation.

Each of these category of transitions are presented in further detail below.

### Housekeeping Transitions:

| Name        | Params     | Description | Callable when paused?|
| ----------- | -----------|-------------|:--------------------------:|
| `Pause` | `initiator : ByStr20`| Pause the contract temporarily to stop any critical transition from being invoked. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: | 
| `Unpause` | `initiator : ByStr20`| Un-pause the contract to re-allow the invocation of all transitions. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: |
| `UpdateAdmin` | `admin : ByStr20, initiator : ByStr20` | Replace the current `contractadmin` by `admin`. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: |
| `UpdateVerifier` | `verif : ByStr20, initiator : ByStr20` | Replace the current `verifier` by `verif`. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: |
| `UpdateVerifierRewardAddr` | `addr : ByStr20, initiator : ByStr20` | Replace the current reward receiving address `verifier_receiving_addr` for the `verifier` by `addr`. Since the verifier is currently run by Zilliqa Research, its receiving address is currently updated by the admin not the verifier itself. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: |
| `UpdateStakingParameters` | `min_stake: ByStr20, min_deleg_stake : Uint128, max_comm_change_rate : Uint128, initiator : ByStr20` | Replace the current values of the fields `minstake`, `mindelegstake`,  and `maxcommchangerate` to the input values.  <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: |
| `ChangeBNumReq` | `input_bnum_req : Uint128, initiator : ByStr20` | Replace the current value of the field `bnum_req`.  <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: |
| `UpdateContractAddr` | `proxy_addr : ByStr20, gzil_addr : ByStr20, initiator : ByStr20` | Replace the address of the proxy contract (`proxyaddr`) and the gZIL token contract (`gziladdr`) by the input values. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: |
| `AddSSN` | `ssnaddr : ByStr20, name : String, urlraw : String, urlapi : String, comm : Uint128, initiator : ByStr20` | Add a new SSN to the list of available SSNs with the input values. The transition will create a value of type `ssn` using the input values and add it to the `ssnlist` map. Since, this SSN is new, the `status` field in `ssn` type will be `False`. Similarly, the fields `stake_amt`, `rewards`, `buff_deposit`, `comm_rewards` in the `ssn` type will be set to 0. The commission for this SSN for this reward cycle as recorded in the field `comm_for_ssn` will also be set to 0. The transition will emit a success event to signal the addition of the new SSN. The event will emit the address of the new SSN. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: |
| `UpdateSSN` | `ssnaddr : ByStr20, new_name : String, new_urlraw : String, new_urlapi : String, initiator : ByStr20` | Update certain fields corresponding to a given SSN.  Other fields not covered by the input parameters should remain unchanged. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: |
| `RemoveSSN` | `ssnaddr : ByStr20, initiator : ByStr20` | Remove a specific SSN from `ssnlist`. It should also remove the corresponding entry from all other fields. **This is a very privileged operation and therefore should not be invoked until all delegators have withdrawn their stake and their stake rewards.  If the SSN gets removed while there exists a delegator at this SSN, then funds for the delegator may get locked forever.**  <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: |

### Delegator Transitions

| Name        | Params     | Description | Callable when paused?|
| ----------- | -----------|-------------|:--------------------------:|
| `DelegateStake` | `ssnaddr : ByStr20, initiator : ByStr20` | To delegate the stake to an SSN. `initiator` is the address of the delegator. The stake being delegated is captured in the implicit field `_amount`. In case of failure to accept the stake, `_amount` must be returned back to the delegator. This stake will be buffered if the SSN is already active else it will be added to the stake pool of the SSN. This transition must update the following fields accordingly: `ssnlist` (with the updated buffered stake or unbuffered stake), `deposit_amt_deleg` and `ssn_deleg_amt` (with the amount deposited by the delegator), `buff_deposit_deleg` (with the amount deposited in case the SSN was active), `direct_deposit_deleg` (with the amount deposited in case the SSN was inactive), `last_withdraw_cycle_deleg` (to record the reward cycle number when the deposit was made for first-time delegators) and `last_buf_deposit_cycle_deleg` (to record the reward cycle number in case the SSN was active and the stake ended up as buffered deposit). The transition should throw an error in case the amount being delegated is less than `mindelegstake`.   | <center>:x:</center> |
| `WithdrawStakeRewards` | `ssnaddr : ByStr20, initiator : ByStr20` | To withdraw stake rewards from a given SSN. `initiator` is the address of the delegator. Reward to be given to a delegator can be computed on-the-fly proportional to its unbuffered stake in each cycle. I.e., if the delegator's unbuffered stake at this SSN in a given cycle is d and the total unbuffered stake across all SSNs is D, and the total reward available for all SSNs is R, then the reward earned by this delegator for the cycle will be `(R * d)/D`. When a delegator calls this transition, the reward is computed for all cycles i such that `lastWithdrawnCycle < i <=  lastrewardcyle`. The transition also mints gZILs and assigns them to the delegator's address. For every ZIL earned as reward, the delegator will earn 0.001 gZIL. | <center>:x:</center> |
| `WithdrawStakeAmt` | `ssnaddr : ByStr20, amt : Uint128, initiator : ByStr20` | To withdraw stake amount from a given SSN. `initiator` is the address of the delegator. No stake amount can be withdrawn if the delegator has unwithdrawn reward or buffered deposit. If the amount withdrawn is less than or equal to the stake deposited, then the withdrawal request is accepted and the field `withdrawal_pending` is updated accordingly. A delegator may request multiple withdrawals in the same cycle. Once this transition is called, the delegator enters into an unbonding period of 24,000 blocks. Only at the expiry of the unbonding period, the delegator can claim the funds. During the unbonding period, the delegator will not earn any reward. | <center>:x:</center> |
| `CompleteWithdrawal` | `initiator : ByStr20` | This is to be called after the delegator has called `WithdrawStakeAmt`. `initiator` is the address of the delegator. The transition processing all the pending withdrawals as recorded in `withdrawal_pending`. Only those pending requests for which the bonding period has expired will be processed. Upon success, all the funds get transferred from the contract to the `initiator`. | <center>:x:</center> |


### Pause Transitions

| Name        | Params     | Description | Callable when paused?|
| ----------- | -----------|-------------|:--------------------------:|
 `update_minstake` | `min_stake : Uint128, initiator : ByStr20` | Update the value of the field `min_stake` to the input value `min_stake`. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| <center>:x:</center> |
| `update_maxstake` | `max_stake : Uint128, initiator : ByStr20` | Update the value of the field `max_stake` to the input value `max_stake`. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| <center>:x:</center> |
| `update_contractmaxstake` | `max_stake : Uint128, initiator : ByStr20` | Update the value of the field `contractmaxstake` to the input value `max_stake`. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| <center>:x:</center> |
| `drain_contract_balance` | `initiator : ByStr20` | Allows the admin to withdraw the entire balance of the contract. It should only be invoked in case of emergency. The withdrawn ZILs go to a multsig wallet contract that represents the `admin`. :warning: **Note:** `initiator` must be the current `contractadmin` of the contract. | :heavy_check_mark:|



### SSN Operation Transitions

| Name        | Params     | Description | Callable when paused?|
| ----------- | -----------|-------------|:--------------------------:|
| `UpdateComm` | `new_rate : Uint128, initiator : ByStr20`| To update the commission rate. `initiator` is the SSN operator. An operator cannot update twice in the same cycle. The `new_rate` must also be less that the field `maxcommrate` and the change in the rate compared from the old one must be less than or equal to `maxcommchangerate`. | <center>:x:</center> | 
| `WithdrawComm` | `ssnaddr : ByStr20, initiator : ByStr20`| To withdraw the commission earned. `initiator` is the SSN's commission receiving address. In case the `initiator` is not the receiving address for the SSN or if the SSN does not exist or if the SSN hasn't earned any commission, then the transition throws error. On success, the contract transfer the commission to the `initiator`. | <center>:x:</center> | 
| `UpdateReceivingAddr` | `newaddr : ByStr20, initiator : ByStr20`| To update the commission receiving address for the SSN. `initiator` is the address of the SSN. | <center>:x:</center> | 


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
