package transitions

import "log"

func (t *Testing) AssignStakeReward() {
	t.LogStart("AssignStakeReward")
	proxy, ssnlist := t.DeployAndUpgrade()
	// unpause
	proxy.Unpause()
	// set staking parameters
	min := "100000000000000"
	max := "500000000000000"
	cmax := "1000000000000000"
	proxy.UpdateStakingParameters(min, max, cmax)
	// update verifier to addr2
	proxy.UpdateVerifier("0x" + addr2)
	// add ssn1
	proxy.AddSSN("0x"+addr1, "ssn1")
	// fund ssnlist
	proxy.AddFunds("100000000000000")
	// delegate stake
	proxy.DelegateStake("0x"+addr1, "100000000000000")
	ssnlist.LogContractStateJson()

	// use addr2(which is verifier) to assign rewards
	proxy.UpdateWallet(key2)
	txn,err := proxy.AssignStakeReward("0x"+addr1, "52000000")
	if err != nil {
		t.LogError("AssignStakeReward",err)
	}
	receipt :=  t.GetReceiptString(txn)
	log.Println(receipt)
	state := ssnlist.LogContractStateJson()
	t.AssertContain(state,"\"lastrewardcycle\":\"1\"")
	t.AssertContain(state,"5200000000000")
	t.AssertContain(state,"\"reward_cycle_list\":[\"1\"]")

	// use addr1 (which is not verifier) to assign rewards
	proxy.UpdateWallet(key1)
	txn,err1 := proxy.AssignStakeReward("0x"+addr1, "52000000")
	t.AssertError(err1)
	receipt =  t.GetReceiptString(txn)
	log.Println(receipt)

	t.LogEnd("AssignStakeReward")
}
