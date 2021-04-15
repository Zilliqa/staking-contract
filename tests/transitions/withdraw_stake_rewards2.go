package transitions

import "log"

// case for inactive => delegate => active => withdraw(one) => withdraw(none) => reward => reward => withdraw(two)
func (t *Testing) WithdrawStakeReward2() {
	t.LogStart("WithdrawStakeReward2")
	proxy, ssnlist := t.DeployAndUpgrade()
	// unpause
	proxy.Unpause()
	// set staking parameters
	min := "100000000000000"
	delegMin := "50000"
	proxy.UpdateStakingParameters(min, delegMin)
	// update verifier to addr1
	proxy.UpdateVerifier("0x" + addr1)
	// update verifier receiving addr to add1
	proxy.UpdateVerifierRewardAddr("0x" + addr1)
	// add ssn1
	proxy.AddSSNAfterUpgrade("0x"+addr1, "0")
	proxy.AddFunds("1000000000000")
	// delegate stake
	proxy.DelegateStake("0x"+addr1, "100000000000000")
	proxy.AssignStakeRewardFixed("0x"+addr1, "52000000")
	ssnlist.LogContractStateJson()

	// withdraw rewards
	txn, err2 := proxy.WithdrawStakeRewards("0x" + addr1)
	if err2 != nil {
		t.LogError("WithdrawStakeReward2", err2)
	}
	receipt := t.GetReceiptString(txn)
	log.Println(receipt)
	ssnlist.LogContractStateJson()

	proxy.WithdrawStakeRewards("0x" + addr1)
	ssnlist.LogContractStateJson()

	proxy.AssignStakeRewardFixed("0x"+addr1, "52000000")
	proxy.AssignStakeRewardFixed("0x"+addr1, "52000000")
	proxy.WithdrawStakeRewards("0x" + addr1)
	ssnlist.LogContractStateJson()

	t.LogEnd("WithdrawStakeReward2")
}
