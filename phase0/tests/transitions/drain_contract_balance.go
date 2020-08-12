package transitions

import (
	"errors"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"strconv"
)

func (p *Proxy) drainContractBalance(pri string) error {
	proxy, _ := bech32.ToBech32Address(p.Addr)
	if err2, _ := ExecZli("contract", "call",
		"-k", pri,
		"-a", proxy,
		"-t", "drain_contract_balance",
		"-f", "true",
		"-r", "[]"); err2 != nil {
		return errors.New("call transition error: " + err2.Error())
	} else {
		m := p.Provider.GetBalance(p.ImplAddress).Result.(map[string]interface{})
		r := m["balance"].(string)
		newbalance, err := strconv.ParseInt(r, 10, 64)
		if err != nil {
			return errors.New("parse balance error: " + err.Error())
		}

		if newbalance != 0 {
			return errors.New("darin contract failed")
		}

		return nil
	}
}
