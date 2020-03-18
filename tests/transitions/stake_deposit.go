//
// test: stake_deposit
//
package transitions

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/util"
)

// preset scenarios before testing stake_deposit
// adjust the values according to values in updateMinStake() and updateMaxStake() calls
const (
	MAX_STAKE = 5000000000
	MIN_STAKE = 1000000000
)

func TestStakeDeposit(pri1, pri2 string, api string) {
	fmt.Println("------------------------ start stakeDeposit ------------------------")
	err, proxy, impl := DeployAndUpgrade(pri1)
	if err != nil {
		panic("got error = " + err.Error())
	}
	fmt.Println("proxy = ", proxy)
	fmt.Println("impl = ", impl)
	p := NewProxy(api, proxy, impl)
	p.StakeDeposit(pri1, pri2)
}

func (p *Proxy) StakeDeposit(pri1, pri2 string) {
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
	err = p.updateVerifier(pri1)
	if err != nil {
		panic("update verifier error: " + err.Error())
	}

	// 1. non-ssn transfer min_stake amount into contract
	proxy, _ := bech32.ToBech32Address(p.Addr)
	if err2, output := ExecZli("contract", "call",
		"-k", pri1,
		"-a", proxy,
		"-t", "stake_deposit",
		"-m", strconv.Itoa(MIN_STAKE),
		"-f", "true",
		"-r", "[]"); err2 != nil {
		panic("call transaction error: " + err2.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		eventLogs := receipt["event_logs"].([]interface{})[0]
		if success {
			events := eventLogs.(map[string]interface{})
			eventName := events["_eventname"].(string)
			if eventName == "SSN doesn't exist" {
				fmt.Println("test stake deposit non-registered ssn succeed")
			} else {
				panic("test stake deposit non-registered ssn failed")
			}
		} else {
			panic("test stake deposit non-registered ssn failed")
		}
	}

	// 2. as ssn1 (pri2), add amount below min stake
	// must perform add_ssn first

	// 2.1. add pri2 as ssn
	p.RegisterSSN(pri1, pri2)

	// 2.2 execute stake deposit as ssn (pri2) with min stake - 1
	if err3, output := ExecZli("contract", "call",
		"-k", pri2,
		"-a", proxy,
		"-t", "stake_deposit",
		"-m", strconv.Itoa(MIN_STAKE-1),
		"-f", "true",
		"-r", "[]"); err3 != nil {
		panic("call transaction error: " + err3.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println(tx)
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		eventLogs := receipt["event_logs"].([]interface{})[0]
		if success {
			events := eventLogs.(map[string]interface{})
			eventName := events["_eventname"].(string)
			if eventName == "SSN stake deposit below min_stake limit" {
				fmt.Println("test stake deposit below min stake limit succeed")
			} else {
				panic("test stake deposit below min stake limit failed")
			}
		} else {
			panic("test stake deposit below min stake limit error")
		}
	}

	// 3. as ssn, stake deposit with max state + 1
	if err3, output := ExecZli("contract", "call",
		"-k", pri2,
		"-a", proxy,
		"-t", "stake_deposit",
		"-m", strconv.Itoa(MAX_STAKE+1),
		"-f", "true",
		"-r", "[]"); err3 != nil {
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
			if eventName == "SSN stake deposit above max_stake limit" {
				fmt.Println("test stake deposit above max stake limit succeed")
			} else {
				panic("test stake deposit above max stake limit failed, tx:" + tx)
			}
		} else {
			panic("test stake deposit above max stake limit error, tx:" + tx)
		}
	}

	// 4. as ssn, first time stake deposit (MIN_STAKE + 1)
	// NO SUCCESS OR FAIL EVENT?
	// if err3, output := ExecZli("contract", "call",
	// 	"-k", pri2,
	// 	"-a", proxy,
	// 	"-t", "stake_deposit",
	// 	"-m", strconv.Itoa(MIN_STAKE+1),
	// 	"-f", "true",
	// 	"-r", "[]"); err3 != nil {
	// 	panic("call transaction error: " + err3.Error())
	// } else {
	// 	tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
	// 	fmt.Println(tx)
	// 	payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
	// 	receipt := payload["receipt"].(map[string]interface{})
	// 	success := receipt["success"].(bool)
	// 	eventLogs := receipt["event_logs"].([]interface{})[0]
	// 	if success {
	// 		events := eventLogs.(map[string]interface{})
	// 		eventName := events["_eventname"].(string)

	// 		if eventName == "SSN updated stake" {
	// 			ssnAddr := events["ssn_address"].(string)
	// 			newStakeAmount := events["new_stake_amount"].(string)
	// 			expectedSSNAddr := "0x" + keytools.GetAddressFromPrivateKey(util.DecodeHex(pri2))
	// 			expectedStakeAmount := strconv.Itoa(MIN_STAKE + 1)

	// 			if ssnAddr != expectedSSNAddr {
	// 				panic("test first time stake deposit failed, tx:" + tx + " , returned ssn: " + ssnAddr + " , expected ssn: " + expectedSSNAddr)
	// 			}
	// 			if newStakeAmount != expectedStakeAmount {
	// 				panic("test first time stake deposit failed, tx:" + tx + " , returned stake amount: " + newStakeAmount + " , expected stake amount: " + expectedStakeAmount)
	// 			}
	// 			fmt.Println("test first time stake deposit succeed")
	// 		} else {
	// 			panic("test first time stake deposit failed, tx:" + tx)
	// 		}
	// 	} else {
	// 		panic("test stake deposit below min stake limit error, tx:" + tx)
	// 	}
	// }

	// 5. as ssn, after first time deposit, deposit max_stake + 1

	// 6. as ssn1, after first time deposit, deposit min_stake + 1

	fmt.Println("------------------------ end StakeDeposit ------------------------")
}

func (p *Proxy) RegisterSSN(pri1, pri2 string) {
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
		fmt.Println(tx)
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if success {
			res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
			sshList, ok := res["ssnlist"]
			if !ok {
				panic("register ssn with verifier failed check state failed: no ssnlist, tx:" + tx)
			}
			ssn, ok := sshList.(map[string]interface{})[ssnaddr]
			if !ok {
				panic("register ssn with verifier failed check state failed: no such ssn, tx:" + tx)
			}
			deposit := ssn.(map[string]interface{})["arguments"].([]interface{})[1].(string)
			if deposit != "0" {
				panic("register ssn with verifier failed check state failed: deposit not equal to zero, tx:" + tx)
			}

			inactive := ssn.(map[string]interface{})["arguments"].([]interface{})[0].(map[string]interface{})["constructor"]
			if inactive != "True" {
				panic("register ssn with verifier failed check state failed: deposit not active, tx:" + tx)
			} else {
				fmt.Println("register ssn with verifier succeed")
			}
		} else {
			panic("register ssn with verifier failed, tx:" + tx)
		}
	}
}
