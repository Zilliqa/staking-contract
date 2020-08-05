package transitions

import (
	"Zilliqa/stake-test/deploy"
	"encoding/json"
	"errors"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"log"
	"strings"
)

const failedLog = "upgradeTo FailedNotAdmin"

func (t *Testing) UpgradeTo()  {
	t.LogStart("UpgradeTo")

	log.Println("start to deploy proxy contract")
	proxy,err := deploy.NewProxy(key1)
	if err != nil {
		t.LogError("deploy proxy error = ",err)
	}
	log.Println("deploy proxy succeed, address = ",proxy.Addr)

	proxy.LogContractStateJson()
	// 1. as admin, upgrade it
	fakeImpl := "0x4226fe2ccbe7c67ae25a0b91414f9129a7892bd5"
	args := []core.ContractValue{
		{
			"newImplementation",
			"ByStr20",
			fakeImpl,
		},
	}

	_,err1 := proxy.Call("UpgradeTo",args)
	if err1 != nil {
		t.LogError("UpgradeTo failed",err1)
	}
	proxy.LogContractStateJson()

	// 2. as non-admin, upgrade it
	proxy.UpdateWallet(key2)
	tnx,_ := proxy.Call("UpgradeTo",args)
	receipt,_ := json.Marshal(tnx.Receipt)
	recp := string(receipt)
	log.Println(recp)
	if !strings.Contains(recp,failedLog) {
		t.LogError("UpgradeTo",errors.New("event log failed"))
	}
	proxy.LogContractStateJson()
	t.LogEnd("UpgradeTo")
}
