package transitions

import "log"

func (t *Testing) AssignStakeReward2() {
	t.LogStart("AssignStakeReward2")
	proxy, ssnlist := t.DeployAndUpgrade()
	// unpause
	proxy.Unpause()
	// set staking parameters
	min := "100"
	delegMin := "10"
	proxy.UpdateStakingParameters(min, delegMin)
	// update verifier to addr2
	proxy.UpdateVerifier("0x" + addr2)
	// update verifier receiving addr to addr2
	proxy.UpdateVerifierRewardAddr("0x" + addr2)
	// add ssn1
	proxy.AddSSN("0x"+addr1, "ssn1")
	// add ssn2
	proxy.AddSSN("0x"+addr2, "ssn2")
	// fund ssnlist
	proxy.AddFunds("10000000000")

	ssn1 := "0x" + addr1
	ssn2 := "0x" + addr2

	// for ssn1, deleg1:10, deleg2:100
	proxy.UpdateWallet(key1)
	proxy.DelegateStake(ssn1, "10")
	proxy.UpdateWallet(key2)
	proxy.DelegateStake(ssn1, "100")

	// for ssn2, deleg1:10, deleg2:100, deleg3:10
	proxy.UpdateWallet(key1)
	proxy.DelegateStake(ssn2, "20")
	proxy.UpdateWallet(key2)
	proxy.DelegateStake(ssn2, "100")
	ssnlist.LogContractStateJson()

	proxy.UpdateWallet(key2)
	// reward ssn1
	txn, err := proxy.AssignStakeReward3(ssn1, "100000", ssn2, "100000")
	if err != nil {
		t.LogError("AssignStakeReward2", err)
	}
	receipt := t.GetReceiptString(txn)
	log.Println(receipt)

	// deleg2 withdraw from ssn1 and ssn2
	proxy.UpdateWallet(key2)
	proxy.WithdrawStakeRewards(ssn1)
	proxy.WithdrawStakeRewards(ssn2)

	// as non-verifier
	proxy.UpdateWallet(key1)
	txn, err2 := proxy.AssignStakeRewardFixed("0x"+addr1, "10000000")
	t.AssertError(err2)
	receipt = t.GetReceiptString(txn)
	log.Println(receipt)
	ssnlist.LogContractStateJson()
	t.AssertContain(receipt, "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -2))])")
	t.LogEnd("AssignStakeReward2")
}
