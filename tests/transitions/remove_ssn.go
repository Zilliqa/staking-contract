package transitions

import "log"

func (t *Testing) RemoveSSN() {
	t.LogStart("RemoveSSN")
	proxy, ssnlist := t.DeployAndUpgrade()
	proxy.Unpause()

	// add ssn1,ssn2
	ssnlist.LogContractStateJson()
	proxy.AddSSN("0x"+addr1, "ssn1")
	proxy.AddSSN("0x"+addr2, "ssn2")
	ssnlist.LogContractStateJson()

	// as admin, remove ssn1
	txn, err := proxy.RemoveSSN("0x" + addr1)
	if err != nil {
		t.LogError("RemoveSSN", err)
	}
	receipt := t.GetReceiptString(txn)
	log.Println(receipt)
	t.AssertContain(receipt, "_eventname\":\"SSN removed")
	ssnlist.LogContractStateJson()

	// as admin, remove ssn1 again
	txn, err1 := proxy.RemoveSSN("0x" + addr1)
	t.AssertError(err1)
	receipt = t.GetReceiptString(txn)
	log.Println(receipt)
	t.AssertContain(receipt,"Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -10))])")
	ssnlist.LogContractStateJson()

	// as non admin, remove ssn2
	proxy.UpdateWallet(key2)
	txn, err2 := proxy.RemoveSSN("0x" + addr2)
	t.AssertError(err2)
	receipt = t.GetReceiptString(txn)
	log.Println(receipt)
	t.AssertContain(receipt,"Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -3))])")
	ssnlist.LogContractStateJson()

	t.LogEnd("RemoveSSN")
}
