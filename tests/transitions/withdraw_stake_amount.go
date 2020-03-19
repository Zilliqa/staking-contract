package transitions

import (
	"encoding/json"
	"errors"
	"fmt"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"strconv"
	"strings"
)

func TestWithdrawAmount(pri1, pri2, api string) {
	fmt.Println("------------------------ start TestWithdrawAamount ------------------------")
	// 0. setup
	err, proxy, impl := DeployAndUpgrade(pri1)
	if err != nil {
		panic("got error = " + err.Error())
	}
	fmt.Println("proxy = ", proxy)
	fmt.Println("impl = ", impl)
	p := NewProxy(api, proxy, impl)

	minstake := "100000000000"

	err = p.updateMinStake(pri1, minstake)
	if err != nil {
		panic("update min stake failed: " + err.Error())
	}

	err = p.updateMaxStake(pri1, "20000000000000000000")
	if err != nil {
		panic("update max stake failed: " + err.Error())
	}

	err = p.updateContractMaxStake(pri1, "700000000000000000000")
	if err != nil {
		panic("update contract max stake failed: " + err.Error())
	}

	// 1. no such ssn
	err, event := p.withdrawAmount(pri1, "10000")
	fmt.Println(err.Error())
	fmt.Println(event)
	if event == "SSN doesn't exist" {
		fmt.Println("test withdraw amount (no such ssn) succeed")
	} else {
		panic("test withdraw amount (no such ssn) failed: event error")
	}

	fmt.Println("------------------------ end TestWithdrawAamount ------------------------")

}

func (p *Proxy) withdrawAmount(operator string, amount string) (error, string) {
	res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
	balance := res["_balance"].(string)
	old, _ := strconv.ParseInt(balance, 10, 64)
	proxy := p.Addr
	parameters := []contract2.Value{
		{
			VName: "amount",
			Type:  "Uint128",
			Value: amount,
		},
	}
	args, _ := json.Marshal(parameters)
	if err2, output := ExecZli("contract", "call",
		"-k", operator,
		"-a", proxy,
		"-t", "withdraw_stake_amount",
		"-f", "true",
		"-r", string(args)); err2 != nil {
		return err2, ""
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if success {
			var ename string
			tl := receipt["event_logs"]
			if tl == nil {
				ename = ""
			} else {
				eventLogs := receipt["event_logs"].([]interface{})[0].(map[string]interface{})
				eventName := eventLogs["_eventname"].(string)
				ename = eventName
			}
			res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
			j, _ := json.Marshal(res)
			fmt.Println(string(j))
			balance := res["_balance"].(string)
			newBalance, _ := strconv.ParseInt(balance, 10, 64)
			m := newBalance - old
			expected := strconv.FormatInt(m, 10)
			if expected == amount {
				return nil, ename
			} else {
				return errors.New("check state failed"), ename
			}
		} else {
			return errors.New("transaction failed"), ""
		}

	}
}
