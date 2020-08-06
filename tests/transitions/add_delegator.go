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

	// as admin, add delegator
	tnx, err := proxy.AddDelegator("0x"+addr1, "0xd90f2e538ce0df89c8273cad3b63ec44a3c4ed82", "100000")
	if err != nil {
		t.LogError("AddDelegator", err)
	}
	receipt := t.GetReceiptString(tnx)
	log.Println(receipt)
	state := ssnlist.LogContractStateJson()
	t.AssertContain(state,"\"delegs\":{\"0x29cf16563fac1ad1596dfe6f333978fece9706ec\":{\"0xd90f2e538ce0df89c8273cad3b63ec44a3c4ed82\":\"100000\"}")
	ssnlist.LogContractStateJson()

	// as non admin, add delegator
	proxy.UpdateWallet(key2)
	tnx, err1 := proxy.AddDelegator("0x"+addr1, "0xd90f2e538ce0df89c8273cad3b63ec44a3c4ed83", "100000")
	t.AssertError(err1)
	receipt = t.GetReceiptString(tnx)
	log.Println(receipt)
	t.AssertContain(receipt,"Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -3))])")
	state = ssnlist.LogContractStateJson()
	t.LogEnd("AddDelegator")

}
