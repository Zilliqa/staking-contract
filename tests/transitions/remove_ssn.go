//
// test: remove_deposit
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


func TestRemoveSSN(pri1, pri2 string, api string) {
	fmt.Println("------------------------ start remove ssn ------------------------")
	err, proxy, impl := DeployAndUpgrade(pri1)
	if err != nil {
		panic("got error = " + err.Error())
	}
	fmt.Println("proxy = ", proxy)
	fmt.Println("impl = ", impl)
	p := NewProxy(api, proxy, impl)
	p.RemoveSSN(pri1, pri2, api)
}

func (p *Proxy) RemoveSSN(pri1, pri2 string, api string) {
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


	// 1. add pri2 as ssn
	ssnaddr := "0x" + keytools.GetAddressFromPrivateKey(util.DecodeHex(pri2))
	p.RegisterSSN(pri1, ssnaddr)

	proxy, _ := bech32.ToBech32Address(p.Addr)

	// 2. as ssn, first time stake deposit (MIN_STAKE + 1)
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

	// 3. after first time deposit, deposit min_stake + 1
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

	res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
	r, _ := json.Marshal(res)
	fmt.Println("before removing: ",string(r))


    // 4.6 remove ssn
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
	args, _ := json.Marshal(parameters)
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
			r, _ = json.Marshal(res)
			fmt.Println(string(r))
			ssnList, _ := res["ssnlist"]
			ssnMap := ssnList.(map[string]interface{})
			ssn := ssnMap[ssnaddr]
			totalstakedeposit := res["totalstakedeposit"].(string)
			if ssn == nil && totalstakedeposit == "0" {
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
	
	fmt.Println("------------------------ end remove ssn ------------------------")

}