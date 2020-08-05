package transitions

import (
	"encoding/json"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"log"
)

func (t *Testing) UpdateStakingParameters() {
	t.LogStart("UpdateStakingParameters")
	proxy,ssnlist := t.DeployAndUpgrade()

	// as admin, should succeed
	ssnlist.LogContractStateJson()
	args := []core.ContractValue{
		{
			"min_stake",
			"Uint128",
			"100000",
		},
		{
			"max_stake",
			"Uint128",
			"500000",
		},
		{
			"contract_max_stake",
			"Uint128",
			"1000000",
		},
	}

	txn, err1 := proxy.Call("UpdateStakingParameters", args,"0")
	if err1 != nil {
		t.LogError("UpdateStakingParameters failed", err1)
	}
	receipt, _ := json.Marshal(txn.Receipt)
	recp := string(receipt)
	log.Println(recp)
	state := ssnlist.LogContractStateJson()
	t.AssertContain(state,"\"contractmaxstake\":\"1000000\"")
	t.AssertContain(state,"\"maxstake\":\"500000\"")
	t.AssertContain(state,"\"minstake\":\"100000\"")


	// as non admin
	proxy.UpdateWallet(key2)
	txn, err2 := proxy.Call("UpdateStakingParameters", args,"0")
	t.AssertError(err2)
	receipt, _ = json.Marshal(txn.Receipt)
	recp = string(receipt)
	log.Println(recp)
	ssnlist.LogContractStateJson()

	t.LogEnd("UpdateStakingParameters")
}