package transitions

func (t *Testing) DelegateStake() {
	t.LogStart("DelegateStake")

	// deploy smart contract
	proxy, ssnlist := t.DeployAndUpgrade()
	ssnlist.LogContractStateJson()
	// unpause
	proxy.Unpause()

	// set staking parameters

	t.LogEnd("DelegateStake")

}
