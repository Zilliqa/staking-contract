package transitions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"strconv"
	"strings"
)

// pri1 => admin
// pri2 => operator
func TestWithdrawStakeRewards(pri1, pri2, api string) {
	fmt.Println("------------------------ start WithdrawStakeRewards ------------------------")
	// 0. setup
	err, proxy, impl := DeployAndUpgrade(pri1)
	if err != nil {
		panic("got error = " + err.Error())
	}
	fmt.Println("proxy = ", proxy)
	fmt.Println("impl = ", impl)
	p := NewProxy(api, proxy, impl)

	minstake := "100000000000"

	err = p.updateStakingParameter(pri1, minstake, "20000000000000000000", "700000000000000000000")

	if err != nil {
		panic("update staking parameter error: " + err.Error())
	}

	if err0 := p.unpause(pri1); err0 != nil {
		panic("unpause with valid account failed")
	}

	// 1. no such ssn
	err1, exception := p.withdrawRewards(pri2)
	if err1 != nil && strings.Contains(exception,"SSN doesn't exist") {
		fmt.Println("test withdraw stake rewards (no such ssn) succeed")
	} else {
		panic("test withdraw stake rewards (no such ssn) failed")
	}

	// 2
	// 2.1 update verifier
	err2 := p.updateVerifier(pri1)
	if err2 != nil {
		panic("test withdraw stake rewards failed: update verifier error: " + err2.Error())
	}

	// 2.2 add ssn
	ssnaddr := "0x" + keytools.GetAddressFromPrivateKey(util.DecodeHex(pri2))
	parameters := []contract2.Value{
		{
			VName: "ssnaddr",
			Type:  "ByStr20",
			Value: ssnaddr,
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
	}
	args, _ := json.Marshal(parameters)
	if err2, output := ExecZli("contract", "call",
		"-k", pri1,
		"-a", proxy,
		"-t", "add_ssn",
		"-f", "true",
		"-r", string(args)); err2 != nil {
		panic("test withdraw stake rewards failed: call transaction error: " + err2.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if success {
			fmt.Println("test withdraw stake rewards succeed: add ssn succeed")
		} else {
			panic("test withdraw stake rewards failed: add ssn failed")
		}
	}

	// 2.3 stake deposit
	m := p.Provider.GetBalance(p.ImplAddress).Result.(map[string]interface{})
	r := m["balance"].(string)
	old, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		panic("test withdraw stake rewards failed: parse balance error: " + err.Error())
	}

	if err2, output := ExecZli("contract", "call",
		"-k", pri2,
		"-a", proxy,
		"-t", "stake_deposit",
		"-f", "true",
		"-m", minstake,
		"-r", "[]"); err2 != nil {
		panic("test withdraw stake rewards failed: call transaction error: " + err2.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
		m := p.Provider.GetBalance(p.ImplAddress).Result.(map[string]interface{})
		r := m["balance"].(string)
		newbalance, err := strconv.ParseInt(r, 10, 64)
		if err != nil {
			panic("test withdraw stake rewards failed: parse balance error: " + err.Error())
		}
		delta := newbalance - old
		d := strconv.FormatInt(delta, 10)
		if d != minstake {
			panic("test withdraw stake rewards failed: check state failed")
		} else {
			fmt.Println("test withdraw stake rewards succeed: stake deposit succeed")
		}
	}

	// 2.4 deposit fund(rewards)
	err = p.transferFunds(pri1, "5000")
	if err != nil {
		panic("deposit funds failed: " + err.Error())
	}

	// 2.5 reward
	err = p.assignStakeReward(pri1, ssnaddr, "50")
	if err != nil {
		panic("reward failed: " + err.Error())
	}

	// 2.6 withdraw 100000000000 and rewards
	parameters = []contract2.Value{
		{
			VName: "amount",
			Type:  "Uint128",
			Value: "100000000000",
		},
	}
	args, _ = json.Marshal(parameters)
	if err2, output := ExecZli("contract", "call",
		"-k", pri2,
		"-a", proxy,
		"-t", "withdraw_stake_amount",
		"-f", "true",
		"-r", string(args)); err2 != nil {
		panic("test withdraw stake rewards failed: call transaction error: " + err2.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
		res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
		j, _ := json.Marshal(res)
		fmt.Println(string(j))
		balance := res["_balance"].(string)
		ssnmap := res["ssnlist"].(map[string]interface{})
		ssn := ssnmap[ssnaddr]
		if ssn != nil && balance == "5000" {
			fmt.Println("test withdraw all succeed")
		} else {
			panic("test withdraw all failed: state check failed")
		}
	}

	if err2, output := ExecZli("contract", "call",
		"-k", pri2,
		"-a", proxy,
		"-t", "withdraw_stake_rewards",
		"-f", "true",
		"-r", "[]"); err2 != nil {
		panic("test withdraw stake rewards failed: withdraw rewards error: " + err2.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
		res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
		j, _ := json.Marshal(res)
		fmt.Println(string(j))
		ssnmap := res["ssnlist"].(map[string]interface{})
		ssn := ssnmap[ssnaddr]
		if ssn == nil {
			fmt.Println("test withdraw stake rewards succeed")
		} else {
			panic("test withdraw stake rewards failed: check ssn list state failed")
		}
	}

	parameters2 := []contract2.Value{
		{
			VName: "ssnaddr",
			Type:  "ByStr20",
			Value: ssnaddr,
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
	}
	args2, _ := json.Marshal(parameters2)

	// 2.7 add ssn
	if err2, output := ExecZli("contract", "call",
		"-k", pri1,
		"-a", proxy,
		"-t", "add_ssn",
		"-f", "true",
		"-r", string(args2)); err2 != nil {
		panic("test withdraw stake rewards failed: call transaction error: " + err2.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if success {
			fmt.Println("test withdraw stake rewards succeed: add ssn succeed")
		} else {
			panic("test withdraw stake rewards failed: add ssn failed")
		}
	}

	// 2.8 deposit stake
	m2 := p.Provider.GetBalance(p.ImplAddress).Result.(map[string]interface{})
	r2 := m2["balance"].(string)
	old, _ = strconv.ParseInt(r2, 10, 64)
	if err2, output := ExecZli("contract", "call",
		"-k", pri2,
		"-a", proxy,
		"-t", "stake_deposit",
		"-f", "true",
		"-m", minstake,
		"-r", "[]"); err2 != nil {
		panic("test withdraw stake rewards failed: call transaction error: " + err2.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
		m := p.Provider.GetBalance(p.ImplAddress).Result.(map[string]interface{})
		r := m["balance"].(string)
		newbalance, err := strconv.ParseInt(r, 10, 64)
		if err != nil {
			panic("test withdraw stake rewards failed: parse balance error: " + err.Error())
		}
		delta := newbalance - old
		d := strconv.FormatInt(delta, 10)
		if d != minstake {
			panic("test withdraw stake rewards failed: check state failed")
		} else {
			fmt.Println("test withdraw stake rewards succeed")
		}
	}

	err = p.transferFunds(pri1, "50000")
	if err != nil {
		panic("deposit funds failed: " + err.Error())
	}

	err = p.assignStakeReward(pri1, ssnaddr, "50")
	if err != nil {
		panic("reward failed: " + err.Error())
	}

	// 2.9 withdraw rewards
	if err2, output := ExecZli("contract", "call",
		"-k", pri2,
		"-a", proxy,
		"-t", "withdraw_stake_rewards",
		"-f", "true",
		"-r", "[]"); err2 != nil {
		panic("test withdraw stake rewards failed: withdraw rewards error: " + err2.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
		res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
		j, _ := json.Marshal(res)
		fmt.Println(string(j))
		ssnmap := res["ssnlist"].(map[string]interface{})
		ssn := ssnmap[ssnaddr]
		if ssn == nil {
			panic("test withdraw stake rewards failed: check ssn list state failed")
		} else {
			arguments := ssn.(map[string]interface{})["arguments"].([]interface{})[2].(string)
			if arguments == "0" {
				fmt.Println("test withdraw stake rewards succeed")
			} else {
				panic("test withdraw stake rewards failed: check ssn rewards failed")
			}
		}
	}

}

func (p *Proxy) withdrawRewards(operator string) (error, string) {
	proxy, _ := bech32.ToBech32Address(p.Addr)
	if err2, output := ExecZli("contract", "call",
		"-k", operator,
		"-a", proxy,
		"-t", "withdraw_stake_rewards",
		"-f", "true",
		"-r", "[]"); err2 != nil {
		return errors.New("call transition error: " + err2.Error()), ""
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if success {
			eventLogs := receipt["event_logs"].([]interface{})[0].(map[string]interface{})
			eventName := eventLogs["_eventname"].(string)
			return nil, eventName
		} else {
			exceptions := receipt["exceptions"]
			j, _ := json.Marshal(exceptions)
			return errors.New("transaction failed"), string(j)
		}
	}
}
