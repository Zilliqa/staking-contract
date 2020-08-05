package transitions

import (
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"log"
)

func (t *Testing) Pause() {
	t.LogStart("Pause")

	proxy,ssnlist := t.DeployAndUpgrade()

	// as non admin, pause
	ssnlist.LogContractStateJson()
	proxy.UpdateWallet(key2)
	args := []core.ContractValue{}
	tnx,err := proxy.Call("Pause",args,"0")
	t.AssertError(err)
	receipt :=  t.GetReceiptString(tnx)
	log.Println(receipt)
	t.AssertContain(receipt,"Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -3))])")
	ssnlist.LogContractStateJson()

	// as admin, pause
	proxy.UpdateWallet(key1)
	tnx,err2 := proxy.Call("Pause",args,"0")
	if err2 != nil {
		t.LogError("Pause",err2)
	}
	receipt = t.GetReceiptString(tnx)
	log.Println(receipt)
	ssnlist.LogContractStateJson()

	t.LogEnd("Pause")


}
