package transitions

import "log"

func (t *Testing) DelegateStake() {
	t.LogStart("DelegateStake")

	// deploy smart contract
	proxy, ssnlist := t.DeployAndUpgrade()
	ssnlist.LogContractStateJson()
	// unpause
	proxy.Unpause()

	// set staking parameters
	min := "100000000000000"
	proxy.UpdateStakingParameters(min)

	// add ssn
	proxy.AddSSN("0x"+addr2,"xiaohuo")
	ssnlist.LogContractStateJson()

	// use addr1 to deposit (should enter direct deposit map)
	// ssn becomes active
	txn,err := proxy.DelegateStake("0x"+addr2,"100000000000000")
	if err != nil {
		t.LogError("DelegateStake",err)
	}
	receipt :=  t.GetReceiptString(txn)
	log.Println(receipt)
	state := ssnlist.LogContractStateJson()
	t.AssertContain(state,"_balance\":\"100000000000000")
	t.AssertContain(state,"delegs\":{\"0xe2cd74983c7a3487af3a133a3bf4e7dd76f5d928\":{\"0x29cf16563fac1ad1596dfe6f333978fece9706ec\":\"100000000000000\"}")
	t.AssertContain(state,"deposit_amt_deleg\":{\"0x29cf16563fac1ad1596dfe6f333978fece9706ec\":{\"0xe2cd74983c7a3487af3a133a3bf4e7dd76f5d928\":\"100000000000000\"}")
	t.AssertContain(state,"direct_deposit_deleg\":{\"0x29cf16563fac1ad1596dfe6f333978fece9706ec\":{\"0xe2cd74983c7a3487af3a133a3bf4e7dd76f5d928\":{\"1\":\"100000000000000\"}}")

	// use addr1 to deposit again (should enter buffer deposit map)
	txn,err2 := proxy.DelegateStake("0x"+addr2,"100000000000000")
	if err2 != nil {
		t.LogError("DelegateStake",err2)
	}
	receipt =  t.GetReceiptString(txn)
	log.Println(receipt)
	state = ssnlist.LogContractStateJson()
	t.AssertContain(state,"_balance\":\"200000000000000")
	t.AssertContain(state,"delegs\":{\"0xe2cd74983c7a3487af3a133a3bf4e7dd76f5d928\":{\"0x29cf16563fac1ad1596dfe6f333978fece9706ec\":\"200000000000000\"}")
	t.AssertContain(state,"direct_deposit_deleg\":{\"0x29cf16563fac1ad1596dfe6f333978fece9706ec\":{\"0xe2cd74983c7a3487af3a133a3bf4e7dd76f5d928\":{\"1\":\"100000000000000\"}}")
	t.AssertContain(state,"buff_deposit_deleg\":{\"0x29cf16563fac1ad1596dfe6f333978fece9706ec\":{\"0xe2cd74983c7a3487af3a133a3bf4e7dd76f5d928\":{\"1\":\"100000000000000\"}}")
	t.LogEnd("DelegateStake")

	// delegate to a non-existent ssn, should raise exception
	txn,err3 := proxy.DelegateStake("0x"+addr1,"100000000000000")
	t.AssertError(err3)
	receipt =  t.GetReceiptString(txn)
	log.Println(receipt)
	state = ssnlist.LogContractStateJson()
	t.LogEnd("DelegateStake")
}
