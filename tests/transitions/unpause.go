package transitions

import (
	"log"
)

func (t *Testing) Unpause() {
	t.LogStart("Unpause")

	proxy,ssnlist := t.DeployAndUpgrade()

	// as non admin, unpasue
	ssnlist.LogContractStateJson()
	proxy.UpdateWallet(key2)
	tnx,err := proxy.Unpause()
	t.AssertError(err)
	receipt :=  t.GetReceiptString(tnx)
	log.Println(receipt)
	t.AssertContain(receipt,"Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -3))])")
	ssnlist.LogContractStateJson()

	// as admin, unpause
	proxy.UpdateWallet(key1)
	tnx,err2 := proxy.Unpause()
	if err2 != nil {
		t.LogError("Unpause",err2)
	}
	receipt = t.GetReceiptString(tnx)
	log.Println(receipt)
	ssnlist.LogContractStateJson()
}
