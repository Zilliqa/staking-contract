package transitions

import "log"

func (t *Testing) AddDelegator() {

	t.LogStart("AddDelegator")

	proxy, ssnlist := t.DeployAndUpgrade()

	// add ssn1,ssn2
	ssnlist.LogContractStateJson()
	proxy.AddSSN("0x"+addr1, "ssn1")
	proxy.AddSSN("0x"+addr2, "ssn2")
	ssnlist.LogContractStateJson()

	// as admin, update delegator (addr1) with 100000 to ssn1 (addr1)
	tnx, err := proxy.AddDelegator("0x"+addr1, "0xd90f2e538ce0df89c8273cad3b63ec44a3c4ed82", "100000")
	if err != nil {
		t.LogError("AddDelegator", err)
	}
	receipt := t.GetReceiptString(tnx)
	log.Println(receipt)
	state := ssnlist.LogContractStateJson()
	t.AssertContain(state,"\"deposit_amt_deleg\":{\"0xd90f2e538ce0df89c8273cad3b63ec44a3c4ed82\":{\"0x29cf16563fac1ad1596dfe6f333978fece9706ec\":\"100000\"}")

	// as admin, update delegator (addr1) with 200000 to ssn1 (add1)
	tnx, err1 := proxy.AddDelegator("0x"+addr1, "0xd90f2e538ce0df89c8273cad3b63ec44a3c4ed82", "200000")
	if err1 != nil {
		t.LogError("AddDelegator", err1)
	}
	receipt = t.GetReceiptString(tnx)
	log.Println(receipt)
	state = ssnlist.LogContractStateJson()
	t.AssertContain(state,"\"deposit_amt_deleg\":{\"0xd90f2e538ce0df89c8273cad3b63ec44a3c4ed82\":{\"0x29cf16563fac1ad1596dfe6f333978fece9706ec\":\"300000\"}")

	// as admin, update delegator (addr1) with 0 to ssn1 (add1)
	tnx, err2 := proxy.AddDelegator("0x"+addr1, "0xd90f2e538ce0df89c8273cad3b63ec44a3c4ed82", "0")
	if err2 != nil {
		t.LogError("AddDelegator", err2)
	}
	receipt = t.GetReceiptString(tnx)
	log.Println(receipt)
	t.AssertContain(receipt,"Deleg deleted")
	state = ssnlist.LogContractStateJson()

	// as admin, update delegator (addr2) with 100000 to ssn2 (addr2)
	tnx, err3 := proxy.AddDelegator("0x"+addr2, "0xd90f2e538ce0df89c8273cad3b63ec44a3c4ed81", "100000")
	if err3 != nil {
		t.LogError("AddDelegator", err3)
	}
	receipt = t.GetReceiptString(tnx)
	t.AssertContain(receipt,"Deleg added")
	log.Println(receipt)
	state = ssnlist.LogContractStateJson()

	// as non admin, add delegator
	proxy.UpdateWallet(key2)
	tnx, err4 := proxy.AddDelegator("0x"+addr1, "0xd90f2e538ce0df89c8273cad3b63ec44a3c4ed83", "100000")
	t.AssertError(err4)
	receipt = t.GetReceiptString(tnx)
	log.Println(receipt)
	t.AssertContain(receipt,"Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -3))])")
	state = ssnlist.LogContractStateJson()
	t.LogEnd("AddDelegator")

}
