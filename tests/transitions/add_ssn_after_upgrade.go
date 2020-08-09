package transitions

import "log"

func (t *Testing) AddSSNAfterUpgrade() {
	t.LogStart("AddSSNAfterUpgrade")
	// deploy
	proxy, ssnlist := t.DeployAndUpgrade()


	// unpause
	proxy.Unpause()
	// set staking parameters
	min := "100000000000000"
	more_than_min := "110000000000000"
	ten_zil := "10000000000000"
	proxy.UpdateStakingParameters(min)
	ssnlist.LogContractStateJson()

	// add ssn1 with stake deposit more than min stake
	txn,err := proxy.AddSSNAfterUpgrade("0x"+addr1,more_than_min)
	if err != nil {
		t.LogError("AddSSNAfterUpgrade",err)
	}
	receipt :=  t.GetReceiptString(txn)
	log.Println(receipt)
	state := ssnlist.LogContractStateJson()
	t.AssertContain(state,"\"0x29cf16563fac1ad1596dfe6f333978fece9706ec\":{\"argtypes\":[],\"arguments\":[{\"argtypes\":[],\"arguments\":[],\"constructor\":\"True\"},\"110000000000000\",\"0\",\"fakename\",\"fakeurl\",\"fakeapi\",\"0\",\"0\",\"0\",\"0x29cf16563fac1ad1596dfe6f333978fece9706ec\"],\"constructor\":\"Ssn\"}")

	// add ssn2 with stake deposit less than min stake
	txn,err1 := proxy.AddSSNAfterUpgrade("0x"+addr2,ten_zil)
	if err1 != nil {
		t.LogError("AddSSNAfterUpgrade",err1)
	}
	receipt =  t.GetReceiptString(txn)
	log.Println(receipt)
	state = ssnlist.LogContractStateJson()
	t.AssertContain(state,"0xe2cd74983c7a3487af3a133a3bf4e7dd76f5d928\":{\"argtypes\":[],\"arguments\":[{\"argtypes\":[],\"arguments\":[],\"constructor\":\"False\"},\"10000000000000\",\"0\",\"fakename\",\"fakeurl\",\"fakeapi\",\"0\",\"0\",\"0\",\"0xe2cd74983c7a3487af3a133a3bf4e7dd76f5d928\"],\"constructor\":\"Ssn\"}")

	// add ssn2 again
	txn,err2 := proxy.AddSSNAfterUpgrade("0x"+addr2,ten_zil)
	if err2 != nil {
		t.LogError("AddSSNAfterUpgrade",err2)
	}
	receipt =  t.GetReceiptString(txn)
	log.Println(receipt)
	t.AssertContain(receipt,"_eventname\":\"SSN already exists")
	ssnlist.LogContractStateJson()

	t.LogEnd("AddSSNAfterUpgrade")




}
