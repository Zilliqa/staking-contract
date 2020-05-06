package transitions

import (
	"encoding/json"
	"errors"
	"fmt"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/util"
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

	min := 100000000000
	half := 50000000000
	minstake := strconv.FormatInt(int64(min), 10)
	halfstake := strconv.FormatInt(int64(half), 10)
	therestake := strconv.FormatInt(int64(min*3), 10)
	onehalfstake := strconv.FormatInt(int64(min+half), 10)

	err = p.updateStakingParameter(pri1, minstake, "20000000000000000000", "700000000000000000000")
	if err != nil {
		panic("update staking parameter error: " + err.Error())
	}

	err2 := p.updateVerifier(pri1)
	if err2 != nil {
		panic("test withdraw amount failed: update verifier error: " + err2.Error())
	}

	if err0 := p.unpause(pri1); err0 != nil {
		panic("unpause with valid account failed")
	}

	// 1. no such ssn
	err, exception := p.withdrawAmount(pri2, minstake)
	if err != nil && strings.Contains(exception, "SSN doesn't exist") {
		fmt.Println("test withdraw amount (no such ssn) succeed")
	} else {
		panic("test withdraw amount (no such ssn) failed: event error")
	}

	ssnaddr := "0x" + keytools.GetAddressFromPrivateKey(util.DecodeHex(pri2))
	parameters := []contract2.Value{
		{
			VName: "ssnaddr",
			Type:  "ByStr20",
			Value: ssnaddr,
		},
		{
			VName: "stake_amount",
			Type:  "Uint128",
			Value: "0",
		},
		{
			VName: "rewards",
			Type:  "Uint128",
			Value: "0",
		},
		{
			VName: "urlraw",
			Type:  "String",
			Value: "devapiziiliqacom",
		},
		{
			VName: "urlapi",
			Type:  "String",
			Value: "ziiliqacom",
		},
		{
			VName: "buffered_deposit",
			Type:  "Uint128",
			Value: "0",
		},
	}
	args, _ := json.Marshal(parameters)
	if err2, output := ExecZli("contract", "call",
		"-k", pri1,
		"-a", p.Addr,
		"-t", "add_ssn_after_upgrade",
		"-f", "true",
		"-r", string(args)); err2 != nil {
		panic("call transaction error: " + err2.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if success {
			fmt.Println("test withdraw amount: add ssn succeed")
		} else {
			panic("test withdraw amount failed: add ssn failed")
		}
	}

	// 2 stake deposit
	m := p.Provider.GetBalance(p.ImplAddress).Result.(map[string]interface{})
	r := m["balance"].(string)
	old, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		panic("test withdraw amount failed: parse balance error: " + err.Error())
	}

	if err2, output := ExecZli("contract", "call",
		"-k", pri2,
		"-a", proxy,
		"-t", "stake_deposit",
		"-f", "true",
		"-m", minstake,
		"-r", "[]"); err2 != nil {
		panic("test withdraw amount failed: call transaction error: " + err2.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
		m := p.Provider.GetBalance(p.ImplAddress).Result.(map[string]interface{})
		r := m["balance"].(string)
		newbalance, err := strconv.ParseInt(r, 10, 64)
		if err != nil {
			panic("test withdraw amount failed: parse balance error: " + err.Error())
		}
		delta := newbalance - old
		d := strconv.FormatInt(delta, 10)
		if d != minstake {
			panic("test withdraw amount failed: check state failed")
		} else {
			fmt.Println("test withdraw amount succeed: stake deposit succeed")
		}
	}

	// 3. deposit again to gain buffered deposit
	m = p.Provider.GetBalance(p.ImplAddress).Result.(map[string]interface{})
	r = m["balance"].(string)
	old, err = strconv.ParseInt(r, 10, 64)
	if err2, output := ExecZli("contract", "call",
		"-k", pri2,
		"-a", proxy,
		"-t", "stake_deposit",
		"-f", "true",
		"-m", minstake,
		"-r", "[]"); err2 != nil {
		panic("test withdraw amount failed: call transaction error: " + err2.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
		m := p.Provider.GetBalance(p.ImplAddress).Result.(map[string]interface{})
		r := m["balance"].(string)
		newbalance, err := strconv.ParseInt(r, 10, 64)
		if err != nil {
			panic("test withdraw amount failed: parse balance error: " + err.Error())
		}
		delta := newbalance - old
		d := strconv.FormatInt(delta, 10)
		if d != minstake {
			panic("test withdraw amount failed: check state failed")
		} else {
			fmt.Println("test withdraw amount succeed: stake deposit succeed")
		}
	}

	// 4. withdraw
	err3, _ := p.withdrawAmount(pri2, minstake)
	if err3 != nil {
		fmt.Println("test withdraw amount (no such ssn) succeed")
	} else {
		panic("test withdraw amount (no such ssn) failed: event error")
	}

	// 5 use assign reward to make buffered deposit complete
	err4 := p.assignStakeReward(pri1, ssnaddr, "50")
	if err4 != nil {
		panic("assign reward failed")
	}
	m = p.Provider.GetBalance(p.ImplAddress).Result.(map[string]interface{})
	j, _ := json.Marshal(m)
	fmt.Println(string(j))

	// deposit rewards
	err = p.transferFunds(pri1, "5000")
	if err != nil {
		panic("transfer funds failed")
	}

	// 6 withdraw reward
	err, event := p.withdrawRewards(pri2)
	if err != nil || event != "SSN withdraw reward" {
		fmt.Println("err: " + err.Error())
		fmt.Println("event: " + event)
		panic("withdraw reward error")
	} else {
		res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
		ssnmap := res["ssnlist"].(map[string]interface{})
		ssn := ssnmap[ssnaddr]
		arguments := ssn.(map[string]interface{})["arguments"].([]interface{})[2].(string)
		if arguments == "0" {
			fmt.Println("withdraw reward succeed")
		} else {
			panic("withdraw reward error: rewards should be 0")
		}
	}

	// 7 withdraw 3 * min stake: emit above error
	err5, exception := p.withdrawAmount(pri2, therestake)
	if err5 != nil && strings.Contains(exception, "SSN withdrawal above stake") {
		fmt.Println("withdraw amount succeed: should emit above error")
	} else {
		panic("withdraw amount failed: should emit above error")
	}

	// 8 withdraw 1.5 * min stake: emit below error
	err6, exception := p.withdrawAmount(pri2, onehalfstake)
	if err6 != nil && strings.Contains(exception,"SSN withdrawal below min_stake limit") {
		fmt.Println("withdraw amount succeed: should emit below error")
	} else {
		panic("withdraw amount failed: should emit below error")
	}

	// 9 withdraw half min stake: remain 1.5 min stake
	err9, _ := p.withdrawAmount(pri2, halfstake)
	if err9 != nil {
		fmt.Println("err: " + err9.Error())
		panic("withdraw amount failed: should remain 1.5 min stake")
	} else {
		fmt.Println("withdraw amount succeed: should remain 1.5 min stake")
	}

	// 10 withdraw the rest one and half: should be removed
	err10, _ := p.withdrawAmount(pri2, onehalfstake)
	if err10 != nil {
		fmt.Println("err: " + err10.Error())
		//panic("withdraw amount failed: should remain 1.5 min stake")
	} else {
		res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
		ssnmap := res["ssnlist"].(map[string]interface{})
		ssn := ssnmap[ssnaddr]
		if ssn == nil {
			fmt.Println("withdraw amount succeed: should be removed")
		} else {
			panic("withdraw amount failed: should be removed")
		}
	}
	fmt.Println("------------------------ end TestWithdrawAmount ------------------------")
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
		fmt.Println("transaction id = ", tx)
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
			m := old - newBalance
			expected := strconv.FormatInt(m, 10)
			if expected == amount {
				return nil, ename
			} else {
				return errors.New("check state failed"), ename
			}
		} else {
			exceptions := receipt["exceptions"]
			j, _ := json.Marshal(exceptions)
			return errors.New("transaction failed"), string(j)
		}
	}
}
