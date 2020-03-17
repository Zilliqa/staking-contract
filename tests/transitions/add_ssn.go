package transitions

import (
	"encoding/json"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/util"
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
	// 0. setup minstake maxstake contractmaxstake
	err := p.updateContractMaxStake(pri1)
	if err != nil {
		panic("update contract max stake failed: " + err.Error())
	}
	err = p.updateMinStake(pri1)
	if err != nil {
		panic("update min stake failed: " + err.Error())
	}
	err = p.updateMaxStake(pri1)
	if err != nil {
		panic("update max stake failed: " + err.Error())
	}

	// 1. as non-verifier to add ssn, should fail
	proxy, _ := bech32.ToBech32Address(p.Addr)
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
			panic("test add ssn with non-verifier failed")
		} else {
			fmt.Println("test add ssn with non-verifier succeed")
		}
	}

	//// 2. stake 100
	//if err4, output := ExecZli("contract", "call",
	//	"-k", pri2,
	//	"-a", proxy,
	//	"-t", "stake_deposit",
	//	"-m", "100",
	//	"-f", "true",
	//	"-r", string(args)); err4 != nil {
	//	panic("call deposit stake transaction error: " + err4.Error())
	//} else {
	//	tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
	//	payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
	//	receipt := payload["receipt"].(map[string]interface{})
	//	success := receipt["success"].(bool)
	//	if success {
	//		panic("test stake less than min stake failed")
	//	} else {
	//		fmt.Println("test stake less then min stake succeed")
	//	}
	//}

	// 3. as verifier to add ssn
	// 3.1 update verifier
	err = p.updateVerifier(pri2)
	if err != nil {
		panic("update verifier error: " + err.Error())
	}

	// 3.2 add account2 to ssn list
	if err3, output := ExecZli("contract", "call",
		"-k", pri2,
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
				panic("test add ssn with verifier failed check state failed: no ssnlist")
			}
			ssn, ok := sshList.(map[string]interface{})[ssnaddr]
			if !ok {
				panic("test add ssn with verifier failed check state failed: no such ssn")
			}
			deposit := ssn.(map[string]interface{})["arguments"].([]interface{})[1].(string)
			if deposit != "0" {
				panic("test add ssn with verifier failed check state failed: deposit not equal to zero")
			}

			inactive := ssn.(map[string]interface{})["arguments"].([]interface{})[0].(map[string]interface{})["constructor"]
			if inactive != "True" {
				panic("test add ssn with verifier failed check state failed: deposit not active")
			} else {
				fmt.Println("test add ssn with verifier succeed")
			}
		} else {
			panic("test add ssn with verifier failed")
		}

		//// 4. stake 100 (less then min stake, should fail)
		//if err4, output := ExecZli("contract", "call",
		//	"-k", pri2,
		//	"-a", proxy,
		//	"-t", "stake_deposit",
		//	"-m", "100",
		//	"-f", "true",
		//	"-r", string(args)); err4 != nil {
		//	panic("call deposit stake transaction error: " + err4.Error())
		//} else {
		//	tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		//	payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		//	receipt := payload["receipt"].(map[string]interface{})
		//	success := receipt["success"].(bool)
		//	if success {
		//		panic("test stake less than min stake failed")
		//	} else {
		//		fmt.Println("test stake less then min stake succeed")
		//	}
		//}

		// 3.3 add ssn once again
		if err3, output := ExecZli("contract", "call",
			"-k", pri2,
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
				panic("test add ssn twice failed")
			} else {
				fmt.Println("test add ssn twice succeed")
			}
		}

		// 3.4 remove an nonexistent
		parameters := []contract2.Value{
			{
				VName: "ssnaddr",
				Type:  "ByStr20",
				Value: "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			},
		}
		args, _ := json.Marshal(parameters)
		if err3, output := ExecZli("contract", "call",
			"-k", pri2,
			"-a", proxy,
			"-t", "remove_ssn",
			"-f", "true",
			"-r", string(args)); err3 != nil {
			panic("call transaction error: " + err3.Error())
		} else {
			tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
			payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
			receipt := payload["receipt"].(map[string]interface{})
			success := receipt["success"].(bool)
			eventLogs := receipt["event_logs"].([]interface{})[0]
			if success {
				events := eventLogs.(map[string]interface{})
				eventName := events["_eventname"].(string)
				if eventName == "SSN already exists" {
					fmt.Println("test add ssn twice succeed")
				} else {
					panic("test add ssn twice failed")
				}
			} else {
				panic("test add ssn twice failed")
			}
		}

	}

	fmt.Println("------------------------ end   AddSSN ------------------------")
}
