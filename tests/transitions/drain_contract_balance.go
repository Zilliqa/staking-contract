package transitions

import (
	"errors"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"log"
)

func (t Testing) DrainContractBalance() {
	t.LogStart("DrainContractBalance")

	proxy,ssnlist := t.DeployAndUpgrade()

	// unpasue
	proxy.Unpause()

	// fund 100zil first
	proxy.AddFunds("100000000000000")

	// drain 5000
	args := []core.ContractValue{
		{
			"amt",
			"Uint128",
			"50000000000000",
		},
	}

	// drain balance while contract is unpaused
	tnx,err := proxy.Call("DrainContractBalance",args,"0")
	t.AssertError(err)

	// pause contract
	proxy.Pause()

	// drain balance again
	tnx,err1 := proxy.Call("DrainContractBalance",args,"0")

	if err1 != nil {
		t.LogError("DrainContractBalance",err)
	}

	if ssnlist.GetBalance() != "50000000000000" {
		t.LogError("DrainContractBalance",errors.New("balance error"))
	}

	receipt :=  t.GetReceiptString(tnx)
	log.Println(receipt)
	t.LogEnd("DrainContractBalance")
}
