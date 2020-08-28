package transitions

import "log"

func (t *Testing) UpdateReceiveAddr() {
	t.LogStart("UpdateReceiveAddr")
	proxy, ssnlist := t.DeployAndUpgrade()
	proxy.AddSSN("0x"+addr1, "xiaohuo")
	proxy.Unpause()

	ssnlist.LogContractStateJson()

	// as ssn operator, update receiving address
	newAddr := "0x" + addr2
	txn, err := proxy.UpdateReceiveAddr(newAddr)
	if err != nil {
		t.LogError("UpdateReceiveAddr", err)
	}
	receipt := t.GetReceiptString(txn)
	log.Println(receipt)
	state := ssnlist.LogContractStateJson()
	t.AssertContain(state,newAddr)

	// as non operator, should fail
	proxy.UpdateWallet(key2)
	txn, err1 := proxy.UpdateReceiveAddr(newAddr)
	t.AssertError(err1)
	receipt = t.GetReceiptString(txn)
	log.Println(receipt)
	t.AssertContain(receipt,"Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -9))])")

	t.LogEnd("UpdateReceiveAddr")
}
