package transitions

import (
	"Zilliqa/stake-test/deploy"
	"encoding/json"
	"log"

	"github.com/Zilliqa/gozilliqa-sdk/core"
)

const adminChanged = "ChangeProxyAdmin"
const adminNotChanged = "ChangeProxyAdmin FailedNotAdmin"

func (t *Testing) ChangeProxyAdmin() {
	t.LogStart("ChangeProxyAdmin")
	log.Println("start to deploy proxy contract")
	proxy, err := deploy.NewProxy(key1)
	if err != nil {
		t.LogError("deploy proxy error = ", err)
	}
	log.Println("deploy proxy succeed, address = ", proxy.Addr)

	proxy.LogContractStateJson()
	// 1. as admin, change proxy admin
	newAdmin := "0x4226fe2ccbe7c67ae25a0b91414f9129a7892bd5"
	args := []core.ContractValue{
		{
			"newAdmin",
			"ByStr20",
			newAdmin,
		},
	}

	txn, err1 := proxy.Call("ChangeProxyAdmin", args, "0")
	if err1 != nil {
		t.LogError("ChangeProxyAdmin failed", err1)
	}

	receipt, _ := json.Marshal(txn.Receipt)
	recp := string(receipt)
	t.AssertContain(recp, adminChanged)
	log.Println(recp)
	proxy.LogContractStateJson()

	// 2. as non-admin, change proxy admin
	proxy.UpdateWallet(key2)
	tnx, _ := proxy.Call("ChangeProxyAdmin", args, "0")
	receipt, _ = json.Marshal(tnx.Receipt)
	recp = string(receipt)
	log.Println(recp)
	t.AssertContain(recp, adminNotChanged)
	proxy.LogContractStateJson()

	t.LogEnd("ChangeProxyAdmin")
}
