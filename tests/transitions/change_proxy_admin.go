package transitions

import (
	"Zilliqa/stake-test/deploy"
	"encoding/json"
	"log"

	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/util"
)

const adminChanged = "ChangeProxyAdmin"
const proxyAdminClaimed = "ClaimProxyAdmin"
const notStagingAdmin = "ClaimProxyAdmin FailedNotStagingadmin"

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
		t.LogError("ChangeProxyAdmin FailedNotAdmin", err1)
	}

	receipt, _ := json.Marshal(txn.Receipt)
	recp := string(receipt)
	t.AssertContain(recp, adminChanged)
	log.Println(recp)
	proxy.LogContractStateJson()

	// 2. as non-admin, claim proxy admin
	args2 := []core.ContractValue{}
	tnx, _ := proxy.Call("ClaimProxyAdmin", args2, "0")
	receipt, _ = json.Marshal(tnx.Receipt)
	recp = string(receipt)
	log.Println(recp)
	t.AssertContain(recp, notStagingAdmin)
	proxy.LogContractStateJson()

	// 3. change back to original staging admin
	newStagingAdmin := keytools.GetAddressFromPrivateKey(util.DecodeHex(key2))
	args = []core.ContractValue{
		{
			"newAdmin",
			"ByStr20",
			"0x" + newStagingAdmin,
		},
	}
	proxy.Call("ChangeProxyAdmin", args, "0")

	// as key2, claim proxy admin
	proxy.UpdateWallet(key2)
	tnx, _ = proxy.Call("ClaimProxyAdmin", args2, "0")
	receipt, _ = json.Marshal(tnx.Receipt)
	recp = string(receipt)
	log.Println(recp)
	t.AssertContain(recp, proxyAdminClaimed)
	proxy.LogContractStateJson()

	t.LogEnd("ChangeProxyAdmin")
}
