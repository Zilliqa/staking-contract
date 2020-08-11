package transitions

import "log"

func (t *Testing) AssignStakeReward2() {
	t.LogStart("AssignStakeReward2")
	proxy, ssnlist := t.DeployAndUpgrade()
	// unpause
	proxy.Unpause()
	// set staking parameters
	min := "100000000000000"
	proxy.UpdateStakingParameters(min)
	// update verifier to addr2
	proxy.UpdateVerifier("0x" + addr2)
	// add ssn1
	proxy.AddSSN("0x"+addr1, "ssn1")
	// add ssn2
	proxy.AddSSN("0x"+addr2, "ssn2")
	// fund ssnlist
	proxy.AddFunds("100000000000000")
	// delegate stake
	proxy.AddDelegator("0x"+addr1, "0x"+addr3, "100000000000000")
	proxy.AddDelegator("0x"+addr2, "0x"+addr3, "100000000000000")
	ssnlist.LogContractStateJson()

	proxy.UpdateWallet(key2)
	// reward ssn1
	txn, err := proxy.AssignStakeReward("0x"+addr1, "10000000")
	if err != nil {
		t.LogError("AssignStakeReward2", err)
	}
	receipt := t.GetReceiptString(txn)
	log.Println(receipt)
	state := ssnlist.LogContractStateJson()
	t.AssertContain(state,"\"0x29cf16563fac1ad1596dfe6f333978fece9706ec\":{\"argtypes\":[],\"arguments\":[{\"argtypes\":[],\"arguments\":[],\"constructor\":\"True\"},\"100000000000000\",\"1000000000000\",\"ssn1\",")
	t.AssertContain(state,"0xe2cd74983c7a3487af3a133a3bf4e7dd76f5d928\":{\"argtypes\":[],\"arguments\":[{\"argtypes\":[],\"arguments\":[],\"constructor\":\"True\"},\"100000000000000\",\"0\",\"ssn2\"")
	// reward ssn1 and ssn2
	txn, err1 := proxy.AssignStakeReward2("0x"+addr1, "10000000","0x"+addr2,"10000000")
	if err1 != nil {
		t.LogError("AssignStakeReward2",err1)
	}
	receipt = t.GetReceiptString(txn)
	log.Println(receipt)
	state = ssnlist.LogContractStateJson()
	t.AssertContain(state,"\"0x29cf16563fac1ad1596dfe6f333978fece9706ec\":{\"argtypes\":[],\"arguments\":[{\"argtypes\":[],\"arguments\":[],\"constructor\":\"True\"},\"100000000000000\",\"2000000000000\",\"ssn1\"")
	t.AssertContain(state,"\"0xe2cd74983c7a3487af3a133a3bf4e7dd76f5d928\":{\"argtypes\":[],\"arguments\":[{\"argtypes\":[],\"arguments\":[],\"constructor\":\"True\"},\"100000000000000\",\"1000000000000\",\"ssn2\"")
	t.AssertContain(state,"\"lastrewardcycle\":\"2\"")

	// as non-verifier
	proxy.UpdateWallet(key1)
	txn, err2 := proxy.AssignStakeReward("0x"+addr1, "10000000")
	t.AssertError(err2)
	receipt = t.GetReceiptString(txn)
	log.Println(receipt)
	state = ssnlist.LogContractStateJson()
	t.AssertContain(receipt,"Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -2))])")
	t.LogEnd("AssignStakeReward2")
}
