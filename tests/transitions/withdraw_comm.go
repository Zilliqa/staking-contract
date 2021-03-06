package transitions

import (
	"encoding/json"
	"log"
)

func (t *Testing) WithdrawComm() {
	t.LogStart("WithdrawComm")
	proxy,ssnlist := t.DeployAndUpgrade()

	// unpause
	proxy.Unpause()
	// set staking parameters
	min := "100000000000000"
	delegMin := "50000"
	proxy.UpdateStakingParameters(min,delegMin)
	// update verifier to addr2
	proxy.UpdateVerifier("0x" + addr1)
	// update verifier receiving addr to add1
	proxy.UpdateVerifierRewardAddr("0x" + addr1)
	// add ssn1
	proxy.AddSSN("0x"+addr1, "ssn1")
	// pause
	proxy.Pause()
	proxy.PopulateTotalStakeAmt("1000000000000000")
	// unpause
	proxy.Unpause()
	// delegate stake
	proxy.AddDelegator("0x"+addr1, "0x"+addr3, "100000000000000")
	proxy.AssignStakeReward("0x"+addr1, "10000000")
	proxy.UpdateComm("100000000")
	// fund ssnlist
	proxy.AddFunds("100000000000000")
	proxy.AssignStakeReward("0x"+addr1, "10000000")
	ssnlist.LogContractStateJson()

	txn,err := proxy.WithdrawComm("0x"+addr1)
	if err != nil {
		t.LogError("WithdrawComm",err)
	}
	receipt, _ := json.Marshal(txn.Receipt)
	recp := string(receipt)
	log.Println(recp)
	ssnlist.LogContractStateJson()

	proxy.UpdateWallet(key2)
	txn,err1 := proxy.WithdrawComm("0x"+addr1)
	t.AssertError(err1)
	ssnlist.LogContractStateJson()

	t.LogEnd("WithdrawComm")
}