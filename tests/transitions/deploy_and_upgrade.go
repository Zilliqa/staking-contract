package transitions

import (
	"Zilliqa/stake-test/deploy"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"log"
)

// this is a help function : )
func (t *Testing) DeployAndUpgrade() (*deploy.Proxy, *deploy.SSNList) {
	log.Println("start to deploy proxy contract")
	proxy, err := deploy.NewProxy(key1)
	if err != nil {
		t.LogError("deploy proxy error = ", err)
	}
	log.Println("deploy proxy succeed, address = ", proxy.Addr)

	log.Println("start to deploy ssnlist contract")
	ssnlist, err1 := deploy.NewSSNList(key1, proxy.Addr)
	if err1 != nil {
		t.LogError("deploy ssnlist error = ", err1)
	}
	log.Println("deploy ssnlist succeed, address = ", ssnlist.Addr)

	log.Println("start to upgrade")
	args := []core.ContractValue{
		{
			"newImplementation",
			"ByStr20",
			"0x" + ssnlist.Addr,
		},
	}
	_, err2 := proxy.Call("UpgradeTo", args)
	if err2 != nil {
		t.LogError("UpgradeTo failed", err2)
	}
	log.Println("upgrade succeed")

	return proxy, ssnlist
}
