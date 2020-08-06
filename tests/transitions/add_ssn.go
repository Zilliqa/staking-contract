package transitions

import "log"

func (t *Testing) AddSSN() {
	t.LogStart("AddSSN")
	proxy,ssnlist := t.DeployAndUpgrade()

	ssnlist.LogContractStateJson()
	// as admin
	txn,err := proxy.AddSSN("0xd90f2e538ce0df89c8273cad3b63ec44a3c4ed82","xiaohuo")
	if err != nil {
		t.LogError("AddSSN",err)
	}
	receipt :=  t.GetReceiptString(txn)
	log.Println(receipt)
	state := ssnlist.LogContractStateJson()
	t.AssertContain(state,"{\"argtypes\":[],\"arguments\":[],\"constructor\":\"False\"}")
	t.AssertContain(state,"{\"argtypes\":[],\"arguments\":[{\"argtypes\":[],\"arguments\":[],\"constructor\":\"False\"},\"0\",\"0\",\"xiaohuo\",\"fakeurl\",\"fakeapi\",\"0\",\"0\",\"0\",\"0xd90f2e538ce0df89c8273cad3b63ec44a3c4ed82\"")
	t.LogEnd("AddSSN")

	// as admin, add again, should fail
	txn,err1 := proxy.AddSSN("0xd90f2e538ce0df89c8273cad3b63ec44a3c4ed82","xiaohuo")
	t.AssertError(err1)
	receipt =  t.GetReceiptString(txn)
	t.AssertContain(receipt,"Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -12))])")
	log.Println(receipt)
	state = ssnlist.LogContractStateJson()

	// as non admin, add another, should fail
	proxy.UpdateWallet(key2)
	txn,err2 := proxy.AddSSN("0xd90f2e538ce0df89c8273cad3b63ec44a3c4ed81","xiaohuo1")
	t.AssertError(err2)
	receipt =  t.GetReceiptString(txn)
	t.AssertContain(receipt,"Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -3))])")
	log.Println(receipt)
	ssnlist.LogContractStateJson()

	t.LogEnd("AddSSN")
}
