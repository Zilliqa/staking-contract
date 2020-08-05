package transitions

import "errors"

func (t *Testing) AddFunds() {
	t.LogStart("AddFunds")
	proxy, ssnlist := t.DeployAndUpgrade()
	ssnlist.LogContractStateJson()

	proxy.Unpause()
	proxy.AddFunds("10000")
	ssnlist.LogContractStateJson()

	if ssnlist.GetBalance() != "10000" {
		t.LogError("AddFunds",errors.New("balance error"))
	}

	t.LogEnd("AddFunds")
}
