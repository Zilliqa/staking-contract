package transitions

import (
	"encoding/json"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"strings"
)

// since this serial of test is complicated, so we test them on new fresh contracts
// so we deploy and upgrade first
func TestAddSSN(pri1, pri2 string, api string) {
	fmt.Println("------------------------ start AddSSN ------------------------")
	err, proxy, impl := DeployAndUpgrade(pri1)
	if err != nil {
		panic("got error = " + err.Error())
	}
	fmt.Println("proxy = ", proxy)
	fmt.Println("impl = ", impl)
	p := NewProxy(api, proxy, impl)
	p.AddSSN(pri1, pri2)
}

func (p *Proxy) AddSSN(pri1, pri2 string) {
	// 1. as non-verifier, should fail
	proxy, _ := bech32.ToBech32Address(p.Addr)
	ssnaddr := "0xced263257fa2d12ed0d1fad74ac036162cec9876"
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
		"-a", proxy,
		"-t", "add_ssn",
		"-f", "true",
		"-r", string(args)); err2 != nil {
		panic("call transaction error: " + err2.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if success {
			panic("add ssn with non-verifier failed")
		}
	}

	// 2. as verifier
	// 2.1 update verifier
	err := p.updateVerifier(pri1)
	if err != nil {
		panic("update verifier error: " + err.Error())
	}

	// 2.2
	if err3, output := ExecZli("contract", "call",
		"-k", pri1,
		"-a", proxy,
		"-t", "add_ssn",
		"-f", "true",
		"-r", string(args)); err3 != nil {
		panic("call transaction error: " + err3.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if success {
			res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
			sshList, ok := res["ssnlist"]
			if !ok {
				panic("check state failed: no ssnlist")
			}
			lists, err := json.Marshal(sshList)
			if err != nil {
				panic("check list failed: " + err.Error())
			}
			fmt.Println(string(lists))

		} else {
			panic("add ssn with verifier failed")
		}
	}

	fmt.Println("------------------------ end   AddSSN ------------------------")
}
