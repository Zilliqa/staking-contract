package transitions

func (t *Testing) WithDrawStakeAmount2() {
	t.LogStart("WithDrawStakeAmount2")
	// deploy
	proxy, ssnlist := t.DeployAndUpgrade()

	// unpause
	proxy.Unpause()
	// set staking parameters
	min := "100000000000000"
	min2 := "200000000000000"
	//tenzil := "10000000000000"
	ssn1 := "0x"+addr1
	proxy.UpdateStakingParameters(min)
	// update verifier to addr1
	proxy.UpdateVerifier(ssn1)
	// add ssn1
	proxy.AddSSN(ssn1, "ssn1")
	proxy.AddFunds(min)
	ssnlist.LogContractStateJson()
	// add delegator (addr1) to ssn1 (addr1) with min zil, make ssn active
	proxy.UpdateWallet(key1)
	proxy.DelegateStake(ssn1,min)
	// delegate again, comes to buffered deposit
	proxy.DelegateStake(ssn1,min)
	ssnlist.LogContractStateJson()

	// try withdraw, should fail
	txn, err := proxy.WithdrawStakeAmount("0x" + addr1,min2)
	t.AssertError(err)
	t.LogPrettyReceipt(txn)

	// reward
	proxy.AssignStakeReward(ssn1, "52000000")

	// try withdraw, should fail (because of rewards)
	txn, err1 := proxy.WithdrawStakeAmount(ssn1,min2)
	t.AssertError(err1)
	t.LogPrettyReceipt(txn)

	// withdraw rewards
	proxy.WithdrawStakeRewards(ssn1)

	// withdraw amount
	txn, err2 := proxy.WithdrawStakeAmount(ssn1,min2)
	if err2 != nil {
		t.LogError("WithDrawStakeAmount2",err2)
	}
	t.LogPrettyReceipt(txn)
	t.LogEnd("WithDrawStakeAmount2")
}
