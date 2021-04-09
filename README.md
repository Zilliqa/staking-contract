# Non-Custodial Seed Node Staking Contracts
![GitHub Actions][github-actions-badge]

[github-actions-badge]: https://github.com/Zilliqa/staking-contract/workflows/Typecheck%20contracts/badge.svg

> For phase 0 contracts, please refer to [here](phase0).

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
    + [Housekeeping Transitions](#housekeeping-transitions)
    + [Delegator Transitions](#delegator-transitions)
    + [SSN Operation Transitions](#ssn-operation-transitions)
    + [Verifier Operation Transitions](#verifier-operation-transitions)
    + [Contract Upgrade Transitions](#contract-upgrade-transitions)
    + [Other Transitions](#other-transitions)
- [SSNListProxy Contract Specification](#ssnlistproxy-contract-specification)
  * [Roles and Privileges](#roles-and-privileges-1)
  * [Immutable Parameters](#immutable-parameters-1)
  * [Mutable Fields](#mutable-fields-1)
  * [Transitions](#transitions-1)
      + [Housekeeping Transitions](#housekeeping-transitions)
      + [Relay Transitions](#relay-transitions)
- [gZILToken Contract Specification](#gziltoken-contract-specification)
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
|SSNListProxy| [`proxy.scilla`](./contracts/proxy.scilla)  | A proxy contract that sits on top of the SSNList contract. Any call to the `SSNList` contract must come from `SSNListProxy`. This contract facilitates upgradeability of the `SSNList` contract in case a bug is found.|
|Wallet| [`multisig_wallet.scilla`](./contracts/multisig_wallet.scilla)  | A multisig wallet contract tailored to work with the `SSNListproxy` contract. Certain transitions in the `SSNListProxy` contract can only be invoked when k-out-of-n users have agreed to do so. This logic is handled using the `Wallet` contract. |
|gZILToken| [`gzil.scilla`](./contracts/gzil.scilla)  | A [ZRC-2](https://github.com/Zilliqa/ZRC/blob/master/zrcs/zrc-2.md) compliant fungible token contract. gZIL tokens represent governance tokens. They are issued alongside staking rewards whenever a delegator withdraws her staking rewards. |

# SSNList Contract Specification

The `SSNList` contract is the main contract that is central to the entire staking
infrastructure.

## Roles and Privileges

The table below describes the roles and privileges that this contract defines:

| Role | Description & Privileges|                                    
| --------------- | ------------------------------------------------- |
| `ssn`           | A registered SSN that provides the seed node service and gets rewarded for the service. An SSN can only be registered by the  `admin`. |
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

3. SsnStakeRewardShare Data Type:

```ocaml
type SsnStakeRewardShare = 
| SsnStakeRewardShare of ByStr20 Uint128 Uint128
```
```ocmal
(*  SSNAddress        : ByStr20 *)
(*                      Address of the SSN. *)
(*  CycleReward       : Uint128 *)
(*                      Integer representation of reward assigned by the verifier to this SSN for this cycle. *)
(*  TotalStakeAmount  : Uint128 *)
(*                      Total stake amount at a specific cycle.                                               *)
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
  | VerifierNotSet (* Verifier's address is not set in the field *)
  | VerifierRecvAddrNotSet (* Verifier's reward address address is not set in the field *)
  | ReDelegInvalidSSNAddr (* Delegator cannot redelegate to same SSN address *)
```
## Immutable Parameters

The table below lists the parameters that are defined at the contract deployment time and hence cannot be changed later on.

| Name | Type | Description |                                    
| ---------------      | ----------|-                                         |
| `init_admin`         | `ByStr20` | The initial admin of the contract.       |
| `init_proxy_address` | `ByStr20` | The initial address of the `SSNListProxy` contract.  |
| `init_gzil_address`  | `ByStr20` | Address of the `gZILToken` contract.  |

## Mutable Fields

The table below presents the mutable fields of the contract and their initial values. 

| Name        | Type       | Initial Value                           | Description                                        |
| ----------- | --------------------|--------------- | -------------------------------------------------- |
| `ssnlist`   | `Map ByStr20 Ssn` | `Emp ByStr20 Ssn` | Mapping between SSN addresses and the corresponding `Ssn` information. |
| `comm_for_ssn`   | `Map ByStr20 (Map Uint32 Uint128)` | `Emp ByStr20 (Map Uint32 Uint128)` | `Map (SSNAddress -> Map (RewardCycleNum -> Commission))` This Map is not used for computing commission fee for an SSN. Its purpose is to act as a placeholder to prevent SSN operators from changing their commission rate multiple times within one reward cycle. |
| `deposit_amt_deleg`   | `Map ByStr20 (Map ByStr20 Uint128)` | `Emp ByStr20 (Map ByStr20 Uint128)` | `Map (DelegatorAddress -> Map (SSNAddress -> AmoutDelegated))` |
| `ssn_deleg_amt`   | `Map ByStr20 (Map ByStr20 Uint128)` | `Emp ByStr20 (Map ByStr20 Uint128)` | `Map (SSNAddress -> Map (DelegatorAddress -> AmountDelegated))` This map does not affect any of the contract operation. It is introduced so that wallet developers can easily query the deposit amount given by a delegator. |
| `buff_deposit_deleg`   | `Map ByStr20 (Map ByStr20 (Map Uint32 Uint128))` | `Emp ByStr20 (Map ByStr20 (Map Uint32 Uint128))` | `Map (DelegatorAddress -> Map (SSNAddress -> Map (RewardCycleNum -> BufferedStakeAmount)))` |
| `direct_deposit_deleg`   | `Map ByStr20 (Map ByStr20 (Map Uint32 Uint128))` | `Emp ByStr20 (Map ByStr20 (Map Uint32 Uint128))` | `Map (DelegatorAddress -> Map (SSNAddress -> Map (RewardCycleNum -> UnBufferedStakeAmount)))` |
| `last_withdraw_cycle_deleg`   | `Map ByStr20 (Map ByStr20 Uint32))` | `Emp ByStr20 (Map ByStr20 Uint32))` | `Map (DelegatorAddress -> Map (SSNAddress -> RewardCycleWhenLastWithdrawn))`. For a new delegator that has never deposited any stake for this SSN, this field will store the reward cycle number during which the delegator successfully deposited its stake. |
| `last_buf_deposit_cycle_deleg`   | `Map ByStr20 (Map ByStr20  Uint32))` | `Emp ByStr20 (Map ByStr20 Uint32))` | `Map (DelegatorAddress -> Map (SSNAddress -> RewardCycleWhenLastDeposited))` |
| `stake_ssn_per_cycle`   | `Map ByStr20 (Map Uint32  SSNCycleInfo))` | `Emp ByStr20 (Map Uint32 SSNCycleInfo))` | `Map (SSNAddress -> Map (RewardCycleNum -> SSNCycleInfo))`. **Note that this data type contains information that corresponds to the end of the cycle when the verifier has distributed the reward. It could therefore be different from the information stored in `ssnlist` particularly its `stake` field which gets updated in the middle of a cycle.** |
|`withdrawal_pending` | `Map ByStr20 (Map BNum Uint128)` | `Emp ByStr20 (Map BNum Uint128)` | `Map (DelegatorAddress -> (BlockNumberWhenRewardWithdrawalRequested -> Amount ))` |
| `bnum_req`  | `Uint128` | `Uin128 24000`       | Bonding period in terms of number of blocks. Set to be equivalent to 14 days. |
| `verifier`   | `Option ByStr20` | `None {ByStr20}` | The address of the `verifier`. |
| `verifier_receiving_addr`   | `Option ByStr20` | `None {ByStr20}` | The address to receive verifier's rewards. |
| `minstake`  | `Uint128` | `Uin128 10000000000000000000`       | Minimum stake required to activate an SSN (1 mil ZIL expressed in `Qa`, where `1 ZIL = 10^12 Qa`). |
| `mindelegstake`  | `Uint128` | `Uin128 10000000000000`       | Minimum stake for a delegator (10 ZIL expressed in Qa where 1 ZIL = 10^12 Qa). |
| `contractadmin` | `ByStr20` |  `init_admin` | Address of the administrator of this contract. |
| `proxyaddr` | `ByStr20` |  `init_proxy_address` | Address of the proxy contract. |
| `gziladdr` | `ByStr20` |  `init_gzil_address` | Address of the gzil contract. |
|`lastrewardcycle` | `Uint32` | `Uint32 1` | The block number when the last reward was distributed. |
|`paused` | `ByStr20` | `True` | A flag to record the paused status of the contract. Certain transitions in the contract cannot be invoked when the contract is paused. |
| `maxcommchangerate`| `Uint128` | `Uint128 1` | The maximum rate change that an SSN is allowed to make across cycles. Set to 1%. |
| `maxcommrate`| `Uint128` | `Uint128 1000000000` | The maximum commission rate that an SSN can charge. Set to 100% but represented as an integer multiplied by 10^7. This prevents SSN from setting commission beyond 100%. |
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

### Housekeeping Transitions

| Name        | Params     | Description | Callable when paused? | Callable when not paused? | 
| ----------- | -----------|-------------|:--------------------------:|:--------------------------:|
| `Pause` | `initiator : ByStr20`| Pause the contract temporarily to stop any critical transition from being invoked. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: | :heavy_check_mark: |
| `Unpause` | `initiator : ByStr20`| Un-pause the contract to re-allow the invocation of all transitions. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: | :heavy_check_mark: |
| `UpdateAdmin` | `admin : ByStr20, initiator : ByStr20` | Set a new `stagingcontractadmin` by `admin`. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: | :heavy_check_mark: |
| `ClaimAdmin` | ` initiator : ByStr20` | Claim to be new `contract admin`. <br>  :warning: **Note:** `initiator` must be the current `stagingcontractadmin` of the contract.| :heavy_check_mark: | :heavy_check_mark: |
| `UpdateVerifier` | `verif : ByStr20, initiator : ByStr20` | Replace the current `verifier` by `verif`. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: | :heavy_check_mark: |
| `UpdateVerifierRewardAddr` | `addr : ByStr20, initiator : ByStr20` | Replace the current reward receiving address `verifier_receiving_addr` for the `verifier` by `addr`. Since the verifier is currently run by Zilliqa Research, its receiving address is currently updated by the admin not the verifier itself. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: | :heavy_check_mark: |
| `UpdateStakingParameters` | `min_stake: ByStr20, min_deleg_stake : Uint128, max_comm_change_rate : Uint128, initiator : ByStr20` | Replace the current values of the fields `minstake`, `mindelegstake`,  and `maxcommchangerate` to the input values.  <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: | :heavy_check_mark: |
| `ChangeBNumReq` | `input_bnum_req : Uint128, initiator : ByStr20` | Replace the current value of the field `bnum_req`.  <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: | :heavy_check_mark: |
| `UpdateGzilAddr` | `gzil_addr : ByStr20, initiator : ByStr20` | Replace the gZIL token contract (`gziladdr`) by the input values. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: | :heavy_check_mark: |
| `AddSSN` | `ssnaddr : ByStr20, name : String, urlraw : String, urlapi : String, comm : Uint128, initiator : ByStr20` | Add a new SSN to the list of available SSNs with the input values. The transition will create a value of type `Ssn` using the input values and add it to the `ssnlist` map. Since, this SSN is new, the `status` field in `Ssn` type will be `False`. Similarly, the fields `stake_amt`, `rewards`, `buff_deposit`, `comm_rewards` in the `ssn` type will be set to 0. The transition will emit a success event to signal the addition of the new SSN. The event will emit the address of the new SSN. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: | :heavy_check_mark: |
| `UpdateSSN` | `ssnaddr : ByStr20, new_name : String, new_urlraw : String, new_urlapi : String, initiator : ByStr20` | Update `name`, `urlraw` and `urlapi` of a given SSN. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: | :heavy_check_mark: |


### Delegator Transitions

| Name        | Params     | Description | Callable when paused? | Callable when not paused? |
| ----------- | -----------|-------------|:--------------------------:|:--------------------------:|
| `DelegateStake` | `ssnaddr : ByStr20, initiator : ByStr20` | To delegate the stake to an SSN. `initiator` is the address of the delegator. The stake being delegated is captured in the implicit field `_amount`. In case of failure to accept the stake, an `exception` will be thrown and the transaction will be reverted. This stake will be buffered if the SSN is already active else it will be added to the stake pool of the SSN. This transition must update the following fields accordingly: `ssnlist` (with the updated buffered stake or unbuffered stake), `deposit_amt_deleg` and `ssn_deleg_amt` (with the amount deposited by the delegator), `buff_deposit_deleg` (with the amount deposited in case the SSN was active), `direct_deposit_deleg` (with the amount deposited in case the SSN was inactive), `last_withdraw_cycle_deleg` (to record the reward cycle number when the deposit was made for first-time delegators) and `last_buf_deposit_cycle_deleg` (to record the reward cycle number in case the SSN was active and the stake ended up as buffered deposit). The transition should throw an error in case the amount being delegated is less than `mindelegstake`.   | <center>:x:</center> | :heavy_check_mark: |
| `WithdrawStakeRewards` | `ssnaddr : ByStr20, initiator : ByStr20` | To withdraw stake rewards from a given SSN. `initiator` is the address of the delegator. Reward to be given to a delegator can be computed on-the-fly proportional to its unbuffered stake in each cycle. I.e., if the delegator's unbuffered stake at this SSN in a given cycle is d and the total unbuffered stake across all SSNs is D, and the total reward available for all SSNs is R, then the reward earned by this delegator for the cycle will be `(R * d)/D`. When a delegator calls this transition, the reward is computed for all cycles i such that `lastWithdrawnCycle < i <=  lastrewardcycle`. The transition also mints gZILs and assigns them to the delegator's address. For every ZIL earned as reward, the delegator will earn 0.001 gZIL. | <center>:x:</center> | :heavy_check_mark: |
| `WithdrawStakeAmt` | `ssnaddr : ByStr20, amt : Uint128, initiator : ByStr20` | To withdraw stake amount from a given SSN. `initiator` is the address of the delegator. No stake amount can be withdrawn if the delegator has unwithdrawn reward or buffered deposit. If the amount withdrawn is less than or equal to the stake deposited, then the withdrawal request is accepted and the field `withdrawal_pending` is updated accordingly. A delegator may request multiple withdrawals in the same cycle. Once this transition is called, the delegator enters into an unbonding period of 24,000 blocks. Only at the expiry of the unbonding period, the delegator can claim the funds via `CompleteWithdrawal` transition. During the unbonding period, the delegator will not earn any reward. | <center>:x:</center> | :heavy_check_mark: |
| `CompleteWithdrawal` | `initiator : ByStr20` | This is to be called after the delegator has called `WithdrawStakeAmt`. `initiator` is the address of the delegator. The transition processing all the pending withdrawals as recorded in `withdrawal_pending`. Only those pending requests for which the bonding period has expired will be processed. Upon success, all the funds get transferred from the contract to the `initiator`. | <center>:x:</center> | :heavy_check_mark: |
| `ReDelegateStake` | `ssnaddr: ByStr20, to_ssn: ByStr20, amount: Uint128, initiator: ByStr20` | To re-delagate the stake from a SSN to another SSN. `initiator` is the address of the delegator. `ssnaddr` is the original SSN, `to_ssn` is the new SSN the delegator wants to delegate to. The re-delegate amount is specificed by `amount`. | <center>:x:</center> | :heavy_check_mark: |
| `RequestDelegatorSwap` | `initiator: ByStr20, new_deleg_addr: ByStr20` | Creates a request to another delegator to indicate transferring all existing stakes, rewards, etc., to this new delegator. `initiator` is the address of the delegator who wants to transfer his/her stakes. `new_deleg_addr` is the address of the recipient that would be receiving all the staked amount, rewards, pending withdrawals, etc., of the `initiator` (original owner). The `initiator` is allowed to change the recipient by sending the request with another `new_deleg_addr`. The `initiator` can also revoke the request via `RevokeDelegatorSwap`. On the recipient end, the `new_deleg_addr` can either `ConfirmDelegatorSwap` to accept the swap or `RejectDelegatorSwap` to reject the swap. One caveat is that the `initiator` is not allowed to make a request to `new_deleg_addr` if there is an existing request made by `new_deleg_addr` to the `initiator`; i.e., if there exists a `A -> B` request, then `B` cannot make a request to `A` unless `B` accepts or rejects the existing request first. However, `B` can make other swap requests to other delegators. **Change is irreversible once the recipient accepts the swap request, please be cautious of the `new_deleg_addr`.** | <center>:x:</center> | :heavy_check_mark: |
| `RevokeDelegatorSwap` | `initiator: ByStr20` | Revokes a swap request. This is used only by the `initiator` who has made an existing swap request and wishes to cancel it. |  <center>:x:</center> | :heavy_check_mark: |
| `ConfirmDelegatorSwap` | `initiator: ByStr20, requestor: ByStr20` |  Accepts a swap request from a requestor. `initiator` is the new delegator that would be inheriting all the staked amount, withdrawals, rewards from `requestor`. `requestor` is the delegator who has initiated a swap request via `RequestDelegatorSwap`. | <center>:x:</center> | :heavy_check_mark: |
| `RejectDelegatorSwap` | `initiator: ByStr20, requestor: ByStr20` | Rejects a swap request from a requestor. Once rejected, the requestor must create the swap request again. `initiator` is the new delegator that would be inheriting all the staked amount, withdrawals, rewards from `requestor`. `requestor` is the delegator who has initiated a swap request via `RequestDelegatorSwap`. | <center>:x:</center> | :heavy_check_mark: |

### SSN Operation Transitions

| Name        | Params     | Description | Callable when paused?| Callable when not paused? |
| ----------- | -----------|-------------|:--------------------------:|:--------------------------:|
| `UpdateComm` | `new_rate : Uint128, initiator : ByStr20`| To update the commission rate. `initiator` is the SSN operator. An operator cannot update twice in the same cycle. The `new_rate` must also be less that the field `maxcommrate` and the change in the rate compared from the old one must be less than or equal to `maxcommchangerate`. | <center>:x:</center> | :heavy_check_mark: |
| `WithdrawComm` | `initiator : ByStr20`| To withdraw the commission earned. `initiator` is the SSN operator. On success, the contract transfer the commission to the receiving address. | <center>:x:</center> | :heavy_check_mark: |
| `UpdateReceivingAddr` | `new_addr : ByStr20, initiator : ByStr20`| To update the commission receiving address for the SSN. `initiator` is the address of the SSN. | <center>:x:</center> | :heavy_check_mark: |

### Verifier Operation Transitions

| Name        | Params     | Description | Callable when paused? | Callable when not paused? |
| ----------- | -----------|-------------|:--------------------------:|:--------------------------:|
| `AssignStakeReward` | `ssnreward_list : List SsnRewardShare, available_reward : Uint128, initiator : ByStr20`| To assign reward to each SSN for this cycle. `ssnreward_list` contains the reward factor for each SSN. In more precise terms, it contains the value `(floor((NumberOfDSEpochsInCurrentCycle x 110,000 x VerificationPassed)))`. This input is then multiplied by `(floor(TotalStakeAtSSN / TotalStakeAcrossAllSSNs))` to compute the reward earned by each SSN. `initiator` is the verifier. The `available_reward` contains the rewards for all SSN as well as verifier. Post this call, any buffered deposit with any SSN must be converted to unbuffered stake deposit. The commission earned by the SSNs must also get updated.  | <center>:x:</center> | :heavy_check_mark: |

### Contract Upgrade Transitions

| Name        | Params     | Description | Callable when paused? | Callable when not paused? |
| ----------- | -----------|-------------|:--------------------------:|:--------------------------:|
| `AddSSNAfterUpgrade` | `ssnaddr: ByStr20, stake_amt: Uint128, rewards: Uint128, name: String, urlraw: String, urlapi: String, buff_deposit: Uint128,  comm: Uint128, comm_rewards: Uint128, rec_addr: ByStr20, initiator: ByStr20`| To add a new SSN to the contract. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: | :heavy_check_mark: |
| `UpdateDeleg` | `ssnaddr: ByStr20, deleg : ByStr20, stake_amt: Uint128, initiator: ByStr20`| To add or remove a delegator for an SSN. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: | :heavy_check_mark: |
| `PopulateStakeSSNPerCycle` | `ssn_addr: ByStr20, cycle: Uint32, totalAmt: Uint128, rewards: Uint128, initiator: ByStr20`| To populate `stake_ssn_per_cycle` map. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: | <center>:x:</center> |
| `PopulateLastWithdrawCycleForDeleg` | `deleg_addr: ByStr20, ssn_addr: ByStr20, cycle: Uint32, initiator: ByStr20`| To populate `last_withdraw_cycle_deleg` map. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: | <center>:x:</center> |
| `PopulateLastBufDepositCycleDeleg` | `deleg_addr: ByStr20, ssn_addr: ByStr20, cycle: Uint32, initiator: ByStr20`| To populate `last_buf_deposit_cycle_deleg` map. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: | <center>:x:</center> |
| `PopulateBuffDeposit` | `deleg_addr: ByStr20, ssn_addr: ByStr20, cycle: Uint32, amt: Uint128, initiator: ByStr20` | To populate `buff_deposit_deleg` map. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: | <center>:x:</center> |
| `PopulateDepositAmtDeleg` | `deleg_addr: ByStr20, ssn_addr: ByStr20, amt: Uint128, initiator: ByStr20` | To populate `deposit_amt_deleg` and `ssn_deleg_amt` map. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: | <center>:x:</center> |
| `PopulateDelegStakePerCycle` | `deleg_addr: ByStr20, ssn_addr: ByStr20, cycle: Uint32, amt: Uint128, initiator: ByStr20` | To populate `deleg_stake_per_cycle` map. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: | <center>:x:</center> |
| `PopulateLastRewardCycle` | `cycle: Uint32, initiator: ByStr20` | To populate `lastrewardcycle` . <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: | <center>:x:</center> |
| `PopulateDirectDeposit` | `deleg_addr: ByStr20, ssn_addr: ByStr20, cycle: Uint32, amt: Uint128, initiator: ByStr20` | To populate `direct_deposit_deleg` map. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: | <center>:x:</center> |
| `PopulateCommForSSN` | `ssn_addr: ByStr20, cycle: Uint32, comm: Uint128, initiator: ByStr20` | To populate `comm_for_ssn` map. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: | <center>:x:</center> |
| `PopulateTotalStakeAmt` | `amt: Uint128, initiator: ByStr20` | To populate `totalstakeamount` field. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: | <center>:x:</center> |

### Other Transitions

| Name        | Params     | Description | Callable when paused? | Callable when not paused? |
| ----------- | -----------|-------------|:--------------------------:|:--------------------------:|
| `AddFunds` | `initiator : ByStr20`| To add funds to the contract. Anyone should be able to add funds to the contract.  | :heavy_check_mark: | :heavy_check_mark: |

# SSNListProxy Contract Specification

`SSNListProxy` contract is a relay contract that redirects calls to it to the `SSNList` contract.

## Roles and Privileges

The table below describes the roles and privileges that this contract defines:

| Role | Description & Privileges|                                    
| --------------- | ------------------------------------------------- |
| `init_admin`           | The initial admin of the contract which is usually the creator of the contract. `init_admin` is also the initial value of admin. |                                 |
| `admin`    | Current `admin` of the contract initialized to `init_admin`. Certain critical actions can only be performed by the `admin`, e.g., changing the current implementation of the `SSNList` contract. |
|`initiator` | The user who calls the `SSNListProxy` contract that in turn calls the `SSNList` contract. |

## Immutable Parameters

The table below lists the parameters that are defined at the contract deployment time and hence cannot be changed later on.

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
|`UpgradeTo`| `newImplementation : ByStr20` |  Change the current implementation address of the `SSNList` contract. <br> :warning: **Note:** Only the `admin` can invoke this transition|
|`ChangeProxyAdmin`| `newAdmin : ByStr20` |  Change the current `stagingadmin` of the contract. <br> :warning: **Note:** Only the `admin` can invoke this transition.|
|`ClaimProxyAdmin` | `` |  Change the current `admin` of the contract. <br> :warning: **Note:** Only the `stagingadmin` can invoke this transition.|


### Relay Transitions

These transitions are meant to redirect calls to the corresponding `SSNList`
contract. While redirecting, the contract prepares the `initiator` value that
is the address of the caller of the `SSNListProxy` contract. The signature of
transitions in the two contracts is exactly the same expect the added last
parameter `initiator` for the `SSNList` contract.

| Transition signature in the `SSNListProxy` contract  | Target transition in the `SSNList` contract |
|--|--|
|`Pause()` | `Pause(initiator : ByStr20)` |
|`UnPause()` | `UnPause(initiator : ByStr20)` |
|`UpdateAdmin(new_admin: ByStr20)` | `UpdateAdmin(admin: ByStr20, initiator : ByStr20)`|
|`ClaimAdmin()` | `ClaimAdmin(initiator : ByStr20)`|
|`UpdateVerifier(verif : ByStr20)` | `UpdateVerifier (verif : ByStr20, initiator: ByStr20)`|
|`UpdateVerifierRewardAddr(addr: ByStr20)` | `UpdateVerifierRewardAddr(addr: ByStr20, initiator : ByStr20)`|
|`UpdateStakingParameters(min_stake: Uint128, min_deleg_stake: Uint128, max_comm_change_rate: Uint128)` | `UpdateStakingParameters(min_stake: Uint128, min_deleg_stake: Uint128, max_comm_change_rate: Uint128, initiator : ByStr20) `|
|`ChangeBNumReq(input_bnum_req: Uint128)` | ` ChangeBNumReq(input_bnum_req: Uint128, initiator : ByStr20)`|
|`UpdateGzilAddr(gzil_addr: ByStr20)` | `UpdateGzilAddr(gzil_addr: ByStr20, initiator : ByStr20)`|
|`AddSSN(ssnaddr: ByStr20, name: String, urlraw: String, urlapi: String, comm: Uint128)` | `AddSSN(ssnaddr: ByStr20, name: String, urlraw: String, urlapi: String, comm: Uint128, initiator : ByStr20)`|
|`UpdateSSN(ssnaddr: ByStr20, new_name: String, new_urlraw: String, new_urlapi: String)` | `UpdateSSN(ssnaddr: ByStr20, new_name: String, new_urlraw: String, new_urlapi: String, initiator : ByStr20)`|
|`UpdateComm(new_rate: Uint128)` | `UpdateComm(new_rate: Uint128, initiator : ByStr20)`|
|`WithdrawComm()` | `WithdrawComm(initiator : ByStr20)`|
|`UpdateReceivingAddr(new_addr: ByStr20)` | `UpdateReceivingAddr(new_addr: ByStr20, initiator : ByStr20)`|
|`DelegateStake(ssnaddr: ByStr20)` | `DelegateStake(ssnaddr: ByStr20, initiator : ByStr20)`|
|`WithdrawStakeRewards(ssnaddr: ByStr20)` | `WithdrawStakeRewards(ssnaddr: ByStr20, initiator : ByStr20)`|
|`WithdrawStakeAmt(ssnaddr: ByStr20, amt: Uint128)` | `WithdrawStakeAmt(ssnaddr: ByStr20, amt: Uint128, initiator : ByStr20)`|
|`CompleteWithdrawal()` | `CompleteWithdrawal(initiator : ByStr20)`|
|`ReDelegateStake(ssnaddr : ByStr20, to_ssn : ByStr20, amount : Uint128)` | `ReDelegateStake(ssnaddr : ByStr20, to_ssn : ByStr20, amount : Uint128, initiator : ByStr20)`|
| `RequestDelegatorSwap(new_deleg: ByStr20)` | `RequestDelegatorSwap(initiator: ByStr20, new_deleg_addr: ByStr20)` |
| `RevokeDelegatorSwap()` | `RevokeDelegatorSwap(initiator: ByStr20)` |
| `ConfirmDelegatorSwap(requestor: ByStr20)` | `ConfirmDelegatorSwap(initiator: ByStr20, requestor: ByStr20)` |
| `RejectDelegatorSwap(requestor: ByStr20)` | `RejectDelegatorSwap(initiator: ByStr20, requestor: ByStr20)` |
|`AssignStakeReward(ssnreward_list: List SsnRewardShare)` | `AssignStakeReward(ssnreward_list: List SsnRewardShare, initiator : ByStr20)`|
|`AddFunds()` | `AddFunds(initiator : ByStr20)`|
|`AddSSNAfterUpgrade(ssnaddr: ByStr20, stake_amt: Uint128, rewards: Uint128, name: String, urlraw: String, urlapi: String, buff_deposit: Uint128,  comm: Uint128, comm_rewards: Uint128, rec_addr: ByStr20)` | `AddSSNAfterUpgrade(ssnaddr: ByStr20, stake_amt: Uint128, rewards: Uint128, name: String, urlraw: String, urlapi: String, buff_deposit: Uint128,  comm: Uint128, comm_rewards: Uint128, rec_addr: ByStr20, initiator : ByStr20)`|
|`AddSSNAfterUpgrade(ssnaddr: ByStr20, stake_amt: Uint128, rewards: Uint128, name: String, urlraw: String, urlapi: String, buff_deposit: Uint128,  comm: Uint128, comm_rewards: Uint128, rec_addr: ByStr20)` | `AddSSNAfterUpgrade(ssnaddr: ByStr20, stake_amt: Uint128, rewards: Uint128, name: String, urlraw: String, urlapi: String, buff_deposit: Uint128,  comm: Uint128, comm_rewards: Uint128, rec_addr: ByStr20, initiator : ByStr20)`|
|`UpdateDeleg(ssnaddr: ByStr20, deleg: ByStr20, stake_amt: Uint128)` | `UpdateDeleg(ssnaddr: ByStr20, deleg: ByStr20, stake_amt: Uint128, initiator : ByStr20)`|
|`PopulateStakeSSNPerCycle(ssn_addr: ByStr20, cycle: Uint32, totalAmt: Uint128, rewards: Uint128)` | `PopulateStakeSSNPerCycle(ssn_addr: ByStr20, cycle: Uint32, totalAmt: Uint128, rewards: Uint128, initiator : ByStr20)`|
|`PopulateLastWithdrawCycleForDeleg(deleg_addr: ByStr20, ssn_addr: ByStr20, cycle: Uint32)` | `PopulateLastWithdrawCycleForDeleg(deleg_addr : ByStr20, ssn_addr: ByStr20, cycle: Uint32, initiator : ByStr20)`|
|`PopulateLastBufDepositCycleDeleg(deleg_addr: ByStr20, ssn_addr: ByStr20, cycle: Uint32)` | `PopulateLastBufDepositCycleDeleg(deleg_addr : ByStr20, ssn_addr: ByStr20, cycle: Uint32, initiator : ByStr20)`|
|`PopulateBuffDeposit(deleg_addr: ByStr20, ssn_addr: ByStr20, cycle: Uint32, amt: Uint128)` | `PopulateBuffDeposit(deleg_addr: ByStr20, ssn_addr: ByStr20, cycle: Uint32, amt: Uint128, initiator : ByStr20)`|
|`PopulateDirectDeposit(deleg_addr: ByStr20, ssn_addr: ByStr20, cycle: Uint32, amt: Uint128)` | `PopulateDirectDeposit(deleg_addr: ByStr20, ssn_addr: ByStr20, cycle: Uint32, amt: Uint128, initiator : ByStr20)`|
|`PopulateDepositAmtDeleg(deleg_addr: ByStr20, ssn_addr: ByStr20, amt: Uint128)` | `PopulateDepositAmtDeleg(deleg_addr: ByStr20, ssn_addr: ByStr20, amt: Uint128, initiator : ByStr20)`|
|`PopulateDelegStakePerCycle(deleg_addr: ByStr20, ssn_addr: ByStr20, cycle: Uint32, amt: Uint128)` | `PopulateDelegStakePerCycle(deleg_addr: ByStr20, ssn_addr: ByStr20, cycle: Uint32, amt: Uint128, initiator : ByStr20)`|
|`PopulateLastRewardCycle(cycle: Uint32)` | `PopulateLastRewardCycle(cycle: Uint32, initiator : ByStr20)`|
|`PopulateCommForSSN(ssn_addr: ByStr20, cycle: Uint32, comm: Uint128)` | `PopulateCommForSSN(ssn_addr: ByStr20, cycle: Uint32, comm: Uint128, initiator : ByStr20)`|
|`PopulateTotalStakeAmt(amt: Uint128)` | `PopulateTotalStakeAmt(amt: Uint128, initiator : ByStr20)`|
|`DrainContractBalance(amt: Uint128)` | `DrainContractBalance(amt: Uint128, initiator : ByStr20)`|

> Note: Any transition in `SSNList` contract that accepts money will have the 
corresponding transition in `SSNListProxy` accept the amount and then pass it 
during the message call.

# gZILToken Contract Specification

`gZILToken` contract is a
[ZRC-2](https://github.com/Zilliqa/ZRC/blob/master/reference/FungibleToken-Mintable.scilla)
compliant token contract with a few minor modifications. The contract defines
two extra immutable variables `init_minter` (of type `ByStr20`) and `end_block`
(of type `BNum`). The former is the initial minter of the contract allowed to
mint tokens, while the latter stores the block number until which minting is
allowed. It also introduces a mutable field `minter`(of type `ByStr20`)
initialized to `init_minter` and a transition `ChangeMinter(new_minter:
ByStr20)` to update the address of the minter.  Since `gZILToken` won't require
buring, the `Burn` transition from the ZRC-2 specification is removed.

The modified `Mint` transition to be called by the `minter`:

```ocaml
transition Mint(recipient: ByStr20, amount: Uint128)
  current_block <- & BLOCKNUMBER;
  is_minting_over = builtin blt end_block current_block;
  match is_minting_over with
  | True =>
  | False =>
    IsMinter;
    AuthorizedMint recipient amount;
    (* Prevent sending to a contract address that does not support transfers of token *)
    msg_to_recipient = {_tag: "RecipientAcceptMint"; _recipient: recipient; _amount: zero; 
                        minter: _sender; recipient: recipient; amount: amount};
    msgs = one_msg msg_to_recipient;
    send msgs
  end
end
```

The `ChangeMinter`transition to be called by the `minter`:

```ocaml
transition ChangeMinter(new_minter: ByStr20, initiator: ByStr20)
  IsOwner initiator;
  minter := new_minter;
  e = {_eventname: "ChangedMinter"; new_minter: new_minter};
  event e  
end
```

As tokens are rewarded when a delegator claims its staking rewards within the
`SSNList` contract, the `minter` of the `gZILToken` contract will be the
address of the `SSNList` contract. 

# Multi-signature Wallet Contract Specification

This contract has two main roles. First, it holds funds that can be paid out to
arbitrary users, provided that enough people from a pre-defined set of owners
have signed off on the payout.

Second, and more generally, it also represents a group of users that can invoke
a transition in another contract only if enough people in that group have
signed off on it. In the staking context, it represents the `admin` in the
`SSNList` contract. This provides added security for the privileged `admin`
role.

## General Flow

Any transaction request (whether transfer of payments or invocation of a
foreign transition) must be added to the contract before signatures can be
collected. Once enough signatures are collected, the recipient (in case of
payments) and/or any of the owners (in the general case) can ask for the
transaction to be executed.

If an owner changes his mind about a transaction, the signature can be revoked
until the transaction is executed.

This wallet does not allow adding or removing owners, or changing the number of
required signatures. To do any of those, perform the following steps:

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
|`SubmitNativeTransaction`| `recipient : ByStr20, amount : Uint128, tag : String` | Submit a request for transfer of native tokens for future signoffs. |
|`SubmitCustomUpgradeToTransaction`| `calleeContract : ByStr20, newImplementation : ByStr20` | Submit a request to invoke the `UpgradeTo` transition in the `SSNListProxy` contract. |
|`SubmitCustomChangeProxyAdminTransaction`| `calleeContract : ByStr20, newAdmin : ByStr20` | Submit a request to invoke the `ChangeProxyAdmin` transition in the `SSNListProxy` contract. |
|`SubmitCustomClaimProxyAdminTransaction`| `calleeContract : ByStr20` | Submit a request to invoke the `ClaimProxyAdmin` transition in the `SSNListProxy` contract. |
|`SubmitCustomChangeMinterTransaction`| `calleeContract : ByStr20, new_minter : ByStr20` | Submit a request to invoke the `ChangeMinter` transition in the `gZILToken` contract. |
|`SubmitCustomPauseTransaction`| `calleeContract : ByStr20` | Submit a request to invoke the `Pause` transition in the `SSNListProxy` contract. |
|`SubmitCustomUnpauseTransaction`| `calleeContract : ByStr20` | Submit a request to invoke the `UnPause` transition in the `SSNListProxy` contract. |
|`SubmitCustomUpdateAdminTransaction`| `calleeContract : ByStr20, admin : ByStr20` | Submit a request to invoke the `UpdateAdmin` transition in the `SSNListProxy` contract. |
|`SubmitCustomClaimAdminTransaction`| `calleeContract : ByStr20 | Submit a request to invoke the `ClaimAdmin` transition in the `SSNListProxy` contract. |
|`SubmitCustomUpdateVerifierTransaction`| `calleeContract : ByStr20, verif : ByStr20` | Submit a request to invoke the `UpdateVerifier` transition in the `SSNListProxy` contract. |
|`SubmitCustomUpdateVerifierRewardAddrTransaction`| `calleeContract : ByStr20, verif : ByStr20` | Submit a request to invoke the `UpdateVerifierRewardAddr` transition in the `SSNListProxy` contract. |
|`SubmitCustomUpdateStakingParametersTransaction`| `calleeContract : ByStr20, min_stake : Uint128, min_deleg_stake : Uint128, max_comm_change_rate : Uint128` | Submit a request to invoke the `UpdateStakingParameters` transition in the `SSNListProxy` contract. |
|`SubmitCustomChangeBNumReqTransaction`| `calleeContract : ByStr20, input_bnum_req : Uint128` | Submit a request to invoke the `ChangeBNumReq` transition in the `SSNListProxy` contract. |
|`SubmitCustomUpdateGzilAddrTransaction`| `calleeContract : ByStr20, gzil_addr: ByStr20` | Submit a request to invoke the `UpdateGzilAddr` transition in the `SSNListProxy` contract. |
|`SubmitCustomAddSSNTransaction`| `calleeContract : ByStr20, ssnaddr : ByStr20, stake_amount : Uint128, rewards : Uint128, urlraw : String, urlapi : String, buffered_deposit : Uint128` | Submit a request to invoke the `AddSSN` transition in the `SSNListProxy` contract. |
|`SubmitCustomUpdateSSNTransaction`| `calleeContract : ByStr20, ssnaddr : ByStr20, new_name : String, new_urlraw : String, new_urlapi : String` | Submit a request to invoke the `UpdateSSN` transition in the `SSNListProxy` contract. |
|`SubmitCustomAddSSNAfterUpgradeTransaction` | `calleeContract: ByStr20, ssnaddr: ByStr20, stake_amt: Uint128, rewards: Uint128, name: String, urlraw: String, urlapi: String, buff_deposit: Uint128, comm: Uint128, comm_rewards: Uint128, rec_addr: ByStr20` | Submit a request to invoke the `AddSSNAfterUpgrade` transition in the `SSNListProxy` contract. |
|`SubmitUpdateDelegTransaction` | `calleeContract: ByStr20, ssnaddr: ByStr20, deleg: ByStr20, stake_amt: Uint128` | Submit a request to invoke the `UpdateDeleg` transition in the `SSNListProxy` contract. |
|`SubmitPopulateStakeSSNPerCycleTransaction`| `calleeContract : ByStr20, ssn_address : ByStr20, cycle : Uint32, totalAmt: Uint128, rewards: Uint128` | Submit a request to invoke the `PopulateStakeSSNPerCycle` transition in the `SSNListProxy` contract. |
|`SubmitPopulateLastWithdrawCycleForDelegTransaction`| `calleeContract : ByStr20, deleg_address : ByStr20, ssn_address : ByStr20, cycle : Uint32` | Submit a request to invoke the `PopulateLastWithdrawCycleForDeleg` transition in the `SSNListProxy` contract. |
|`SubmitPopulateLastBufDepositCycleDelegTransaction`| `calleeContract : ByStr20, deleg_address : ByStr20, ssn_address : ByStr20, cycle : Uint32` | Submit a request to invoke the `PopulateLastBufDepositCycleDeleg` transition in the `SSNListProxy` contract. |
|`SubmitPopulateBufferedDepositTransaction`| `calleeContract : ByStr20, deleg_address : ByStr20, ssn_address : ByStr20, cycle : Uint32, amount : Uint128` | Submit a request to invoke the `PopulateBufferedDeposit` transition in the `SSNListProxy` contract. |
|`SubmitPopulateDirectDepositTransaction`| `calleeContract : ByStr20, deleg_address : ByStr20, ssn_address : ByStr20, cycle : Uint32, amount : Uint128` | Submit a request to invoke the `PopulateDirectDeposit` transition in the `SSNListProxy` contract. |
|`SubmitPopulateDepositAmtDelegTransaction`| `calleeContract: ByStr20, deleg_addr: ByStr20, ssn_addr: ByStr20, amt: Uint128` | Submit a request to invoke the `PopulateDepositAmt` transition in the `SSNListProxy` contract. |
|`SubmitPopulateDelegStakePerCycleTransaction`| `calleeContract: ByStr20, deleg_addr: ByStr20, ssn_addr: ByStr20, cycle: Uint32, amt: Uint128` | Submit a request to invoke the `PopulateDelegStakePerCycle` transition in the `SSNListProxy` contract. |
|`SubmitPopulateLastRewardCycleTransaction`| `calleeContract: ByStr20, cycle: Uint32` | Submit a request to invoke the `PopulateLastRewardCycle` transition in the `SSNListProxy` contract. |
|`SubmitPopulateCommForSSNTransaction`| `calleeContract : ByStr20, ssn_address : ByStr20, cycle : Uint32, comm : Uint128` | Submit a request to invoke the `PopulateCommForSSN` transition in the `SSNListProxy` contract. |
|`SubmitPopulateTotalStakeAmtTransaction`| `calleeContract : ByStr20, amt : Uint128` | Submit a request to invoke the `PopulateTotalStakeAmt` transition in the `SSNListProxy` contract. |
|`SubmitPopulatePendingWithdrawalTransaction`| `calleeContract : ByStr20, ssn_addr : ByStr20, block_number : BNum, stake : Uint128` | Submit a request to invoke the `PopulatePendingWithdrawal` transition in the `SSNListProxy` contract. |
|`SubmitCustomDrainContractBalanceTransaction`| `calleeContract : ByStr20, amt : Uint128` | Submit a request to invoke the `DrainContractBalance` transition in the `SSNListProxy` contract. |

### Action Transitions

| Name | Params | Description |
|--|--|--|
|`SignTransaction`| `transactionId : Uint32` | Sign off on an existing transaction. |
|`RevokeSignature`| `transactionId : Uint32` | Revoke signature of an existing transaction, if it has not yet been executed. |
|`ExecuteTransaction`| `transactionId : Uint32` | Execute signed-off transaction. |
