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
	txn, err1 := proxy.UpdateStakingParameters(min, delegMin)
	if err1 != nil {
		t.LogError("UpdateComm failed", err1)
	}
	receipt, _ := json.Marshal(txn.Receipt)
	recp := string(receipt)
	log.Println(recp)

	// update verifier to ssn1
	proxy.UpdateVerifier("0x" + addr1)
	// update verifier receiving addr to addr1
	proxy.UpdateVerifierRewardAddr("0x" + addr1)

	// pause
	proxy.Pause()
	proxy.PopulateTotalStakeAmt("400000")
	proxy.PopulateCommForSSN("0x"+addr1, "1", "100000000")
	// unpause
	proxy.Unpause()

	// add ssn1
	proxy.AddSSNAfterUpgrade("0x"+addr1, "200000")
	txn, err2 := proxy.UpdateStakingParameters(min, delegMin)
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
	t.AssertContain(recp, "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -9))])")
	ssnlist.LogContractStateJson()

	// reward to increase cycle
	proxy.AssignStakeRewardFixed("0x"+addr1, "52000000")
	ssnlist.LogContractStateJson()

	// update commission with a legal number
	txn, err6 := proxy.UpdateComm("200000000")
	if err6 != nil {
		t.LogError("UpdateComm", err6)
	}

	// update commission again
	txn, err4 := proxy.UpdateComm("300000000")
	t.AssertError(err4)

	// as non ssn, update commission
	proxy.UpdateWallet(key2)
	txn, err5 := proxy.UpdateComm("150000000")
	t.AssertError(err5)
	receipt, _ = json.Marshal(txn.Receipt)
	recp = string(receipt)
	log.Println(recp)
	t.AssertContain(recp, "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -10))])")
	ssnlist.LogContractStateJson()

	t.LogEnd("UpdateComm")
}
