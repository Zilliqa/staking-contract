package transitions

import (
	"encoding/json"
	"log"
)

func (t *Testing) UpdateComm() {
	t.LogStart("UpdateComm")
	proxy, ssnlist := t.DeployAndUpgrade()

	// unpause
	proxy.Unpause()

	min := "100000"
	delegMin := "50000"
	txn, err1 := proxy.UpdateStakingParameters(min,delegMin)
	if err1 != nil {
		t.LogError("UpdateComm failed", err1)
	}
	receipt, _ := json.Marshal(txn.Receipt)
	recp := string(receipt)
	log.Println(recp)

	// update verifier to ssn1
	proxy.UpdateVerifier("0x" + addr1)
	proxy.PopulateTotalStakeAmt("100")


	// add ssn1
	proxy.AddSSNAfterUpgrade("0x"+addr1, "200000")
	txn, err2 := proxy.UpdateStakingParameters(min,delegMin)
	if err2 != nil {
		t.LogError("UpdateComm failed", err2)
	}
	receipt, _ = json.Marshal(txn.Receipt)
	recp = string(receipt)
	ssnlist.LogContractStateJson()

	// ssn1 update commission within this cycle
	txn, err3 := proxy.UpdateComm("100000000")
	t.AssertError(err3)
	receipt, _ = json.Marshal(txn.Receipt)
	recp = string(receipt)
	log.Println(recp)
	t.AssertContain(recp, "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -8))])")
	ssnlist.LogContractStateJson()

	// reward to increase cycle
	proxy.AssignStakeReward("0x"+addr1, "10")
	ssnlist.LogContractStateJson()

	// update commission with a illegal numner
	txn, err6 := proxy.UpdateComm("400000000")
	t.AssertError(err6)


	// update commission again
	txn, err4 := proxy.UpdateComm("100000000")
	if err4 != nil {
		t.LogError("UpdateComm", err4)
	}
	receipt, _ = json.Marshal(txn.Receipt)
	recp = string(receipt)
	log.Println(recp)
	state := ssnlist.LogContractStateJson()
	t.AssertContain(state,"{\"1\":\"0\",\"2\":\"100000000\"}")

	// as non ssn, update commission
	proxy.UpdateWallet(key2)
	txn, err5 := proxy.UpdateComm("150000000")
	t.AssertError(err5)
	receipt, _ = json.Marshal(txn.Receipt)
	recp = string(receipt)
	log.Println(recp)
	t.AssertContain(recp,"Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -9))])")
	ssnlist.LogContractStateJson()

	t.LogEnd("UpdateComm")
}
