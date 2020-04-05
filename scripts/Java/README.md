# Java sample code for SSN operators
This repository contains some samples illustrated how third parties or organizations integrate with Zilliqa staking contract as SSN operators.

## SSNOperator.java
[SSNOperator.java](./src/main/java/com/zilliqa/staking/SSNOperator.java) cointain the sample codes for SSN operators to interact with the SSN smart contract. 

## public String stakeDeposit(String amount, int attempts, int interval)
This function allows the SSN operator to deposit stake amount into the SSN smart contract. 
```java
@param amount   staking amount
@param attempts attempt times for polling transaction
@param interval interval time in seconds between each polling
```

## public String withdrawStakeAmount(String amount, int attempts, int interval)
This function allows the SSN operator to withdraw stake amount *(excluding reward)* from the SSN smart contract. 
```java
@param amount   withdraw amount
@param attempts attempt times for polling transaction
@param interval interval time in seconds between each polling
```

## public String withdrawStakeRewards(int attempts, int interval)
This function allows the SSN operator to withdraw *all* the stake reward from the SSN smart contract.
```java
* @param attempts attempt times for polling transaction
* @param interval interval time in seconds between each polling
```

## Todo
- [ ] Get staked seed node status
- [ ] Get current stake amount
- [ ] Get current stake reward