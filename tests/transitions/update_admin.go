package transitions

import (
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"log"
)

func (t *Testing) UpdateAdmin() {
	t.LogStart("UpdateAdmin")
	proxy, ssnlist := t.DeployAndUpgrade()

	// as admin, set new admin to addr2
	ssnlist.LogContractStateJson()
	args := []core.ContractValue{{
		"admin",
		"ByStr20",
		"0x" + addr2,
	}}
	tnx, err := proxy.Call("UpdateAdmin", args,"0")
	if err != nil {
		t.LogError("UpdateAdmin", err)
	}
	receipt := t.GetReceiptString(tnx)
	log.Println(receipt)
	t.AssertContain(receipt, addr2)
	ssnlist.LogContractStateJson()

	// as non admin, try to set new admin to addr1
	args = []core.ContractValue{{
		"admin",
		"ByStr20",
		"0x" + addr1,
	}}
	tnx, err2 := proxy.Call("UpdateAdmin", args,"0")
	t.AssertError(err2)

	receipt = t.GetReceiptString(tnx)
	log.Println(receipt)
	t.AssertContain(receipt,"Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -3))]")
	ssnlist.LogContractStateJson()

	t.LogEnd("UpdateAdmin")

}
