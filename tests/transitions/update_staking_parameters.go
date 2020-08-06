package transitions

import (
	"encoding/json"
	"log"
)

func (t *Testing) UpdateStakingParameters() {
	t.LogStart("UpdateStakingParameters")
	proxy,ssnlist := t.DeployAndUpgrade()

	// as admin, should succeed
	ssnlist.LogContractStateJson()
	min := "100000"
	txn,err1 := proxy.UpdateStakingParameters(min)
	if err1 != nil {
		t.LogError("UpdateStakingParameters failed", err1)
	}
	receipt, _ := json.Marshal(txn.Receipt)
	recp := string(receipt)
	log.Println(recp)
	state := ssnlist.LogContractStateJson()
	t.AssertContain(state,"\"minstake\":\"100000\"")


	// as non admin
	proxy.UpdateWallet(key2)
	txn, err2 := proxy.UpdateStakingParameters(min)
	t.AssertError(err2)
	receipt, _ = json.Marshal(txn.Receipt)
	recp = string(receipt)
	log.Println(recp)
	ssnlist.LogContractStateJson()

	t.LogEnd("UpdateStakingParameters")
}