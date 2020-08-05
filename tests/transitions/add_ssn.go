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
	t.LogEnd("AddSSN")
}
