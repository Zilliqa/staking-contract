package transitions

import (
	"errors"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"strconv"
)

func (p *Proxy) TransferFunds(pri, funds string) {
	fmt.Println("------------------------ start transfer funds  ------------------------")
	err := p.transferFunds(pri, funds)
	if err != nil {
		panic("test transfer funds failed")
	} else {
		fmt.Println("test transfer funds succeed")
	}
	fmt.Println("------------------------ end transfer funds ------------------------")

}

func (p *Proxy) transferFunds(pri, funds string) error {
	proxy, _ := bech32.ToBech32Address(p.Addr)
	m := p.Provider.GetBalance(p.ImplAddress).Result.(map[string]interface{})
	r := m["balance"].(string)
	old, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return errors.New("parse balance error: " + err.Error())
	}

	if err2, _ := ExecZli("contract", "call",
		"-k", pri,
		"-a", proxy,
		"-t", "deposit_funds",
		"-m", funds,
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
		d := newbalance - old
		delta := strconv.FormatInt(d, 10)
		if delta != funds {
			return errors.New("check state failed")
		} else {
			return nil
		}

	}

}
