package transitions

import (
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"log"
)

func (t *Testing) UpdateAdmin() {
	t.LogStart("UpdateAdmin")
	proxy, ssnlist := t.DeployAndUpgrade()

	// as admin, propose new admin to addr2
	ssnlist.LogContractStateJson()
	args := []core.ContractValue{{
		"new_admin",
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

	// as origin admin, propose again
	args = []core.ContractValue{{
		"new_admin",
		"ByStr20",
		"0x" + addr2,
	}}
	tnx, err2 := proxy.Call("UpdateAdmin", args,"0")
	if err2 != nil {
		t.LogError("UpdateAdmin", err2)
	}

	// addr2 claim to be admin
	proxy.UpdateWallet(key2)
	tnx, err3 := proxy.ClaimAdmin()
	if err3 != nil {
		t.LogError("ClaimAdmin", err3)
	}

	// addr1 (as non-admin), update admin
	proxy.UpdateWallet(key1)
	tnx, err4 := proxy.Call("UpdateAdmin", args,"0")
	t.AssertError(err4)
	receipt = t.GetReceiptString(tnx)
	log.Println(receipt)
	t.AssertContain(receipt,"Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -3))]")
	ssnlist.LogContractStateJson()

	t.LogEnd("UpdateAdmin")

}
