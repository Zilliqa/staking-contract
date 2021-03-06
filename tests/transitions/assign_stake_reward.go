package transitions

import "log"

func (t *Testing) AssignStakeReward() {
	t.LogStart("AssignStakeReward")
	proxy, ssnlist := t.DeployAndUpgrade()
	// unpause
	proxy.Unpause()
	// set staking parameters
	min := "100000000000000"
	delegMin := "50000"
	proxy.UpdateStakingParameters(min,delegMin)
	// update verifier to addr2
	proxy.UpdateVerifier("0x" + addr2)
	// update verifier receiving addr to addr2
	proxy.UpdateVerifierRewardAddr("0x" + addr2)
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
	t.AssertContain(state,"\"lastrewardcycle\":\"2\"")
	t.AssertContain(state,"52000000")
	// use addr1 (which is not verifier) to assign rewards
	proxy.UpdateWallet(key1)
	txn,err1 := proxy.AssignStakeReward("0x"+addr1, "52000000")
	t.AssertError(err1)
	receipt =  t.GetReceiptString(txn)
	log.Println(receipt)

	t.LogEnd("AssignStakeReward")
}
