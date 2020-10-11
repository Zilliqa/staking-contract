package transitions

import "log"

// case for active => delegate => reward => withdraw(non) => reward => withdraw
func (t *Testing) WithdrawStakeReward3() {
	t.LogStart("WithdrawStakeReward3")
	proxy, ssnlist := t.DeployAndUpgrade()
	// unpause
	proxy.Unpause()
	// set staking parameters
	min := "100000000000000"
	delegMin := "50000"
	proxy.UpdateStakingParameters(min,delegMin)
	// update verifier to addr1
	proxy.UpdateVerifier("0x" + addr1)
	// update verifier receiving addr to add1
	proxy.UpdateVerifierRewardAddr("0x" + addr1)
	// add ssn1
	proxy.AddSSNAfterUpgrade("0x"+addr1,"100000000000000")
	ssnlist.LogContractStateJson()

	// delegate stake
	proxy.DelegateStake("0x"+addr1, "100000000000000")
	proxy.AssignStakeReward("0x"+addr1, "52000000")
	proxy.AddFunds("1000000000000")
	ssnlist.LogContractStateJson()

	// withdraw rewards
	txn,err2 := proxy.WithdrawStakeRewards("0x"+addr1)
	if err2 != nil {
		t.LogError("WithdrawStakeReward3",err2)
	}

	receipt :=  t.GetReceiptString(txn)
	log.Println(receipt)
	ssnlist.LogContractStateJson()

	proxy.AssignStakeReward("0x"+addr1, "52000000")
	txn,err3 := proxy.WithdrawStakeRewards("0x"+addr1)
	if err3 != nil {
		t.LogError("WithdrawStakeReward3",err2)
	}

	receipt = t.GetReceiptString(txn)
	log.Println(receipt)
	ssnlist.LogContractStateJson()


	t.LogEnd("WithdrawStakeReward3")
}