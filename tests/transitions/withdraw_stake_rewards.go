package transitions

import (
	"fmt"
	"log"
)

func (t *Testing) WithdrawStakeReward() {
	t.LogStart("WithdrawStakeReward")
	proxy, ssnlist := t.DeployAndUpgrade()
	// unpause
	proxy.Unpause()
	// set staking parameters
	min := "100000000000000"
	delegMin := "50000"
	proxy.UpdateStakingParameters(min,delegMin)
	// update verifier to addr1
	proxy.UpdateVerifier("0x" + addr1)
	// add ssn1
	proxy.AddSSNAfterUpgrade("0x"+addr1,"200000000000000")
	proxy.AddFunds("1000000000000")
	proxy.AssignStakeReward("0x"+addr1, "52000000")
	// delegate stake
	proxy.DelegateStake("0x"+addr1, "100000000000000")
	proxy.AssignStakeReward("0x"+addr1, "52000000")
	ssnlist.LogContractStateJson()
	proxy.AssignStakeReward("0x"+addr1, "52000000")
	// delegator again
	proxy.DelegateStake("0x"+addr1, "100000000000000")
	proxy.AssignStakeReward("0x"+addr1, "52000000")
	ssnlist.LogContractStateJson()
	// reward 3 more
	result := proxy.AssignStakeRewardBatch("0x"+addr1, "52000000")
	for _,r := range result {
		fmt.Println("to see if error")
		fmt.Println(r.ErrMsg)
	}
	ssnlist.LogContractStateJson()
	// withdraw rewards
	txn,err2 := proxy.WithdrawStakeRewards("0x"+addr1)
	if err2 != nil {
		t.LogError("WithdrawStakeReward",err2)
	}
	receipt :=  t.GetReceiptString(txn)
	log.Println(receipt)
	ssnlist.LogContractStateJson()
	t.LogEnd("WithdrawStakeReward")
}