# Java sample code for SSN operators
This repository contains some sample codes, illustrating how third parties or organizations integrate with Zilliqa staking contract as SSN operators.

## SSNOperator.java
[SSNOperator.java](./src/main/java/com/zilliqa/staking/SSNOperator.java) cointains the sample codes for SSN operators to interact with the SSN smart contract. 

### public String stakeDeposit(String amount, int attempts, int interval)
This function allows the SSN operator to deposit stake amount into the SSN smart contract. 
```java
@param amount   staking amount
@param attempts attempt times for polling transaction
@param interval interval time in seconds between each polling
```

### public String withdrawStakeAmount(String amount, int attempts, int interval)
This function allows the SSN operator to withdraw stake amount *(excluding reward)* from the SSN smart contract. 
```java
@param amount   withdraw amount
@param attempts attempt times for polling transaction
@param interval interval time in seconds between each polling
```

### public String withdrawStakeRewards(int attempts, int interval)
This function allows the SSN operator to withdraw *all* the stake reward from the SSN smart contract.
```java
@param attempts attempt times for polling transaction
@param interval interval time in seconds between each polling
```

### public String getStakeAmount()
This function allows the SSN operator to get stake amount *(excluding reward)* from the SSN smart contract.

### public String getStakeBufferedAmount()
This function allows the SSN operator to get buffered stake amount from the SSN smart contract.

### public String getStakeRewards()
This function allows the SSN operator to get stake reward amount from the SSN smart contract.

### public Boolean getActiveStatus()
This function allows the SSN operator to get the current status (active/inactive) of the SSN. An inactive SSN is usually a result of the existing stake amount less than minstake.

### public Boolean getNodeStatus(String publicApi, int tolerance)
This function allows the SSN operator to check whether is the node is in-sync with the network, up to a certain tolerance.
```java
@param publicApi public api endpoint maintained by Zilliqa, typically https://dev-api.zilliqa.com for community testnet and https://api.zilliqa.com for mainnet
@param tolerance tolerance of the gap of the tx block number
```
