package transitions

func (t *Testing) WithDrawStakeAmount3() {
	t.LogStart("WithDrawStakeAmount3")
	// deploy
	proxy, ssnlist := t.DeployAndUpgrade()

	// unpause
	proxy.Unpause()
	// set staking parameters
	min := "100000000000000"
	min2 := "200000000000000"
	min3 := "300000000000000"
	//tenzil := "10000000000000"
	ssn1 := "0x"+addr1
	delegMin := "50000"
	proxy.UpdateStakingParameters(min,delegMin)
	// update verifier to addr1
	proxy.UpdateVerifier(ssn1)
	// add ssn1
	proxy.AddSSN(ssn1, "ssn1")
	proxy.AddFunds(min2)
	ssnlist.LogContractStateJson()
	// add delegator (addr1) to ssn1 (addr1) with min zil, make ssn active
	proxy.UpdateWallet(key1)
	proxy.DelegateStake(ssn1,min2)
	ssnlist.LogContractStateJson()


	// try withdraw more than delegate
	txn, err := proxy.WithdrawStakeAmount("0x" + addr1,min3)
	t.AssertError(err)
	t.LogPrettyReceipt(txn)
	ssnlist.LogContractStateJson()


	// try withdraw half
	txn, err1 := proxy.WithdrawStakeAmount("0x" + addr1,min)
	if err1 != nil {
		t.LogError("WithDrawStakeAmount3",err)
	}
	t.LogPrettyReceipt(txn)
	ssnlist.LogContractStateJson()
	t.LogEnd("WithDrawStakeAmount3")
}
