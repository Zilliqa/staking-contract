//
// test: stake_deposit
//
package transitions

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/util"

	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
)

// preset scenarios before testing stake_deposit
// adjust the values according to values in updateMinStake() and updateMaxStake() calls
// set the contract max stake to be lesser than max stake on purpose to trigger the above contract max stake event
const (
	CONTRACT_MAX_STAKE = 5000000000
	MAX_STAKE          = 4000000000
	MIN_STAKE          = 1000000000
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
	p.StakeDeposit(pri1, pri2, api)
}

func (p *Proxy) StakeDeposit(pri1, pri2 string, api string) {
	if err0 := p.unpause(pri1); err0 != nil {
		panic("unpause with valid account failed")
	}
	// 0. setup minstake maxstake contractmaxstake
	err := p.updateStakingParameter(pri1, strconv.Itoa(MIN_STAKE), strconv.Itoa(MAX_STAKE), strconv.Itoa(CONTRACT_MAX_STAKE))
	if err != nil {
		panic("update staking parameter error: " + err.Error())
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
		fmt.Println("transaction id = ", tx)
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if !success {
			exceptions := receipt["exceptions"]
			j, _ := json.Marshal(exceptions)
			if !strings.Contains(string(j), "SSN doesn't exist") {
				panic("check exception error")
			}
		} else {
			panic("test stake deposit non-registered ssn failed, tx:" + tx)
		}
	}

	// 2. as ssn1 (pri2), add amount below min stake
	// 2.1. add pri2 as ssn
	ssnaddr := "0x" + keytools.GetAddressFromPrivateKey(util.DecodeHex(pri2))
	p.RegisterSSN(pri1, ssnaddr)

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
		fmt.Println("transaction id = ", tx)
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if !success {
			exceptions := receipt["exceptions"]
			j, _ := json.Marshal(exceptions)
			if !strings.Contains(string(j), "SSN stake deposit below min_stake limit") {
				panic("check exception error")
			}
		} else {
			panic("test stake deposit below min stake limit error, tx:" + tx)
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
		fmt.Println("transaction id = ", tx)
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if !success {
			exceptions := receipt["exceptions"]
			j, _ := json.Marshal(exceptions)
			if !strings.Contains(string(j), "SSN stake deposit above max_stake limit") {
				panic("check exception error")
			}
		} else {
			panic("test stake deposit above max stake limit error, tx:" + tx)
		}
	}

	// 4. as ssn, first time stake deposit (MIN_STAKE + 1)
	if err3, output := ExecZli("contract", "call",
		"-k", pri2,
		"-a", proxy,
		"-t", "stake_deposit",
		"-m", strconv.Itoa(MIN_STAKE+1),
		"-f", "true",
		"-r", "[]"); err3 != nil {
		panic("call transaction error: " + err3.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		eventLogs := receipt["event_logs"].([]interface{})[0]
		if success {
			events := eventLogs.(map[string]interface{})
			eventName := events["_eventname"].(string)

			if eventName == "SSN updated stake" {
				ssnAddr := events["params"].([]interface{})[0].(map[string]interface{})["value"].(string)
				newStakeAmount := events["params"].([]interface{})[1].(map[string]interface{})["value"].(string)
				expectedSSNAddr := "0x" + keytools.GetAddressFromPrivateKey(util.DecodeHex(pri2))
				expectedStakeAmount := strconv.Itoa(MIN_STAKE + 1)

				if ssnAddr != expectedSSNAddr {
					panic("test first time stake deposit failed, tx:" + tx + " , returned ssn: " + ssnAddr + " , expected ssn: " + expectedSSNAddr)
				}
				if newStakeAmount != expectedStakeAmount {
					panic("test first time stake deposit failed, tx:" + tx + " , returned stake amount: " + newStakeAmount + " , expected stake amount: " + expectedStakeAmount)
				}

				// check ssn active status
				res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
				ssnmap := res["ssnlist"].(map[string]interface{})
				ssn := ssnmap[ssnAddr]

				if ssn == nil {
					panic("test first time stake deposit failed, tx:" + tx)
				} else {
					ssnStatus := ssn.(map[string]interface{})["arguments"].([]interface{})[0].(map[string]interface{})["constructor"].(string)
					if ssnStatus == "Active" {
						fmt.Println("test first time stake deposit succeed")
					} else {
						panic("test first time stake deposit failed, tx:" + tx)
					}
				}

			} else {
				panic("test first time stake deposit failed, tx:" + tx)
			}
		} else {
			panic("test stake deposit below min stake limit error, tx:" + tx)
		}
	}

	// 5. as ssn, after first time deposit, deposit contract_max_stake + 1
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
		fmt.Println("transaction id = ", tx)
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if !success {
			exceptions := receipt["exceptions"]
			j, _ := json.Marshal(exceptions)
			if !strings.Contains(string(j), "SSN stake deposit above max_stake limit") {
				panic("check exception error")
			}
		} else {
			panic("test stake deposit (after first time deposit) above max stake limit error, tx:" + tx)
		}
	}

	// 6. as ssn1, after first time deposit, deposit min_stake + 1
	if err3, output := ExecZli("contract", "call",
		"-k", pri2,
		"-a", proxy,
		"-t", "stake_deposit",
		"-m", strconv.Itoa(MIN_STAKE+1),
		"-f", "true",
		"-r", "[]"); err3 != nil {
		panic("call transaction error: " + err3.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if success {
			// no event output after the first time stake deposit
			// check ssn active status
			res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
			totalStakeDeposit := res["totalstakedeposit"].(string)
			expectedTotalStakeDeposit := strconv.Itoa((MIN_STAKE + 1) * 2)
			ssnmap := res["ssnlist"].(map[string]interface{})
			ssnAddr := "0x" + keytools.GetAddressFromPrivateKey(util.DecodeHex(pri2))
			ssn := ssnmap[ssnAddr]

			// this is the second time depositing (min stake + 1)
			if totalStakeDeposit != expectedTotalStakeDeposit {
				panic("test stake deposit (after first time deposit) failed, tx:" + tx + " , total stake deposit:" + totalStakeDeposit + " , expected:" + expectedTotalStakeDeposit)
			}

			if ssn == nil {
				panic("test stake deposit (after first time deposit) error, tx:" + tx)
			} else {
				ssnStatus := ssn.(map[string]interface{})["arguments"].([]interface{})[0].(map[string]interface{})["constructor"].(string)
				if ssnStatus == "Active" {
					fmt.Println("test stake deposit (after first time deposit) succeed")
				} else {
					panic("test stake deposit (after first time deposit) error, tx:" + tx)
				}
			}

		} else {
			panic("test stake deposit (after first time deposit) error, tx:" + tx)
		}
	}

	// 7. as ssn, after second time, deposit CONTRACT_MAX_STAKE
	// invoke CONTRACT_MAX_STAKE
	if err3, output := ExecZli("contract", "call",
		"-k", pri2,
		"-a", proxy,
		"-t", "stake_deposit",
		"-m", strconv.Itoa(CONTRACT_MAX_STAKE),
		"-f", "true",
		"-r", "[]"); err3 != nil {
		panic("call transaction error: " + err3.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if !success {
			exceptions := receipt["exceptions"]
			j, _ := json.Marshal(exceptions)
			if !strings.Contains(string(j), "SSN stake deposit above max_stake limit") {
				panic("check exception error")
			}
		} else {
			panic("test stake deposit above contract max stake limit error, tx:" + tx)
		}
	}

	// 8.  middle contract test
	// 8.1 deploy middle contract
	var middleAddr string
	if err, output := ExecZli("contract", "deploy",
		"-k", pri1,
		"-c", "middle_contract.scilla",
		"-i", "middle_contract.json"); err != nil {
		panic("test stake deposit failed: can not deploy middle contract")
	} else {
		res := strings.Split(output, "contract address =  ")
		res = strings.Split(res[1], "{")
		res = strings.Split(res[0], "track")
		middleAddr = strings.TrimSpace(res[0])
	}
	middleBech32, _ := bech32.ToBech32Address(middleAddr)
	fmt.Printf("middle contract address = %s, bech32 = %s\n", middleAddr, middleBech32)

	// 8.2 register middle contract as ssn
	p.RegisterSSN(pri1, "0x"+middleAddr)

	// 8.3 execute stake deposit with min stake - 1
	parameters := []contract2.Value{
		{
			VName: "proxy_address",
			Type:  "ByStr20",
			Value: "0x" + p.Addr,
		},
	}

	args, _ := json.Marshal(parameters)
	if err3, output := ExecZli("contract", "call",
		"-k", pri1,
		"-a", middleBech32,
		"-t", "stake_deposit",
		"-m", strconv.Itoa(MIN_STAKE-1),
		"-f", "true",
		"-r", string(args)); err3 != nil {
		panic("call transaction error: " + err3.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if !success {
			exceptions := receipt["exceptions"]
			j, _ := json.Marshal(exceptions)
			if !strings.Contains(string(j), "SSN stake deposit below min_stake limit") {
				panic("check exception error")
			}
		} else {
			panic("test stake deposit (middle contract) below min stake limit error, tx:" + tx)
		}
	}

	// 8.4. as ssn, stake deposit with max state + 1
	if err3, output := ExecZli("contract", "call",
		"-k", pri1,
		"-a", middleBech32,
		"-t", "stake_deposit",
		"-m", strconv.Itoa(MAX_STAKE+1),
		"-f", "true",
		"-r", string(args)); err3 != nil {
		panic("call transaction error: " + err3.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if !success {
			exceptions := receipt["exceptions"]
			j, _ := json.Marshal(exceptions)
			if !strings.Contains(string(j), "SSN stake deposit above max_stake limit") {
				panic("check exception error")
			}
		}

		// 8.6 remove ssn
		m := p.Provider.GetBalance(p.ImplAddress).Result.(map[string]interface{})
		before := m["balance"].(string)
		fmt.Println("balance bofore removing ssn: ",before)

		// 4.2 remove exist ssn
		parameters := []contract2.Value{
			{
				VName: "ssnaddr",
				Type:  "ByStr20",
				Value: ssnaddr,
			},
		}
		args, _ = json.Marshal(parameters)
		if err3, output := ExecZli("contract", "call",
			"-k", pri1,
			"-a", proxy,
			"-t", "remove_ssn",
			"-f", "true",
			"-r", string(args)); err3 != nil {
			panic("call transaction error: " + err3.Error())
		} else {
			tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
			fmt.Println("transaction id = ", tx)
			payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
			receipt := payload["receipt"].(map[string]interface{})
			success := receipt["success"].(bool)
			if success {
				res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
				r, _ := json.Marshal(res)
				fmt.Println(string(r))
				ssnList, _ := res["ssnlist"]
				ssnMap := ssnList.(map[string]interface{})
				ssn := ssnMap[ssnaddr]
				if ssn == nil {
					m = p.Provider.GetBalance(p.ImplAddress).Result.(map[string]interface{})
					after := m["balance"].(string)
					if after != "0" {
						panic("check balance failed")
					}
				} else {
					fmt.Println("test remove ssn failed: check state failed")
				}

			} else {
				panic("test remove ssn failed")
			}
		}




	}
	fmt.Println("------------------------ end StakeDeposit ------------------------")
}

func (p *Proxy) RegisterSSN(pri1, ssnaddr string) {
	proxy, _ := bech32.ToBech32Address(p.Addr)
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
		panic("call transaction error: " + err2.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
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
			if inactive == "True" {
				panic("register ssn with verifier failed check state failed: deposit are active, should be inactive tx:" + tx)
			} else {
				fmt.Println("register ssn with verifier succeed")
			}
		} else {
			panic("register ssn with verifier failed, tx:" + tx)
		}
	}
}
