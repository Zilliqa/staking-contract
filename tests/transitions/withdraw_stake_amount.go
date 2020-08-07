package transitions

import "log"

func (t *Testing) WithDrawStakeAmount() {
	t.LogStart("WithDrawStakeAmount")
	// deploy
	proxy, ssnlist := t.DeployAndUpgrade()

	// unpause
	proxy.Unpause()
	// set staking parameters
	min := "100000000000000"
	tenzil := "10000000000000"
	proxy.UpdateStakingParameters(min)
	// update verifier to addr1
	proxy.UpdateVerifier("0x" + addr1)
	// add ssn1
	proxy.AddSSN("0x"+addr1, "ssn1")
	proxy.AddSSN("0x"+addr2, "ssn2")
	// add delegator (addr1) to ssn1 (addr1) with 10 zil
	proxy.AddDelegator("0x"+addr1, "0x"+addr1, tenzil)
	// add delegator (addr2) to ssn1 (addr1) with 10 zil
	proxy.AddDelegator("0x"+addr1, "0x"+addr2, tenzil)
	// add delegator (addr3) to ssn1 (addr1) with min zil
	proxy.AddDelegator("0x"+addr1, "0x"+addr3, min)
	// ssn1 becomes active now
	ssnlist.LogContractStateJson()
	// fund min zil
	proxy.AddFunds(min)

	// delegator (addr2) delegate 10 zil, and it should enter in buffered deposit
	proxy.UpdateWallet(key2)
	proxy.DelegateStake("0x"+addr1, tenzil)
	ssnlist.LogContractStateJson()

	// non delegator(addr4) try to withdraw stake, should fail
	proxy.UpdateWallet(key4)
	txn, err := proxy.WithdrawStakeAmount("0x" + addr1)
	t.AssertError(err)
	receipt :=  t.GetReceiptString(txn)
	log.Println(receipt)
	t.AssertContain(receipt,"Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -5))])")
	ssnlist.LogContractStateJson()

	//delegator (addr1) withdraw from ssn1 (addr1), remain active
	proxy.UpdateWallet(key1)
	txn, err1 := proxy.WithdrawStakeAmount("0x" + addr1)
	if err1 != nil {
		t.LogError("WithDrawStakeAmount",err1)
	}
	receipt =  t.GetReceiptString(txn)
	log.Println(receipt)
	t.AssertContain(receipt,"Deleg withdraw deposit")
	ssnlist.LogContractStateJson()

	// delegator (addr2) withdraw from ssn2 (addr2), should fail
	proxy.UpdateWallet(key2)
	txn, err2 := proxy.WithdrawStakeAmount("0x" + addr2)
	t.AssertError(err2)
	receipt =  t.GetReceiptString(txn)
	log.Println(receipt)
	ssnlist.LogContractStateJson()

	// delegator (addr3) withdraw from ssn1 (addr1), should success, and ssn become inactive
	proxy.UpdateWallet(key3)
	txn, err3 := proxy.WithdrawStakeAmount("0x" + addr1)
	if err3 != nil {
		t.LogError("WithDrawStakeAmount",err3)
	}
	receipt =  t.GetReceiptString(txn)
	log.Println(receipt)
	state := ssnlist.LogContractStateJson()
	t.AssertContain(state,"\"ssnlist\":{\"0x29cf16563fac1ad1596dfe6f333978fece9706ec\":{\"argtypes\":[],\"arguments\":[{\"argtypes\":[],\"arguments\":[],\"constructor\":\"False\"},\"10000000000000\",\"0\",\"ssn1\",\"fakeurl\",\"fakeapi\",\"10000000000000\",\"0\",\"0\",\"0x29cf16563fac1ad1596dfe6f333978fece9706ec\"],\"constructor\":\"Ssn\"}")

	t.LogEnd("WithDrawStakeAmount")
}
