package transitions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"strings"
)

// pri1 admin and verifier
// pri2 aan operator
func TestAssignStakeReward(pri1, pri2, api string) {
	fmt.Println("------------------------ start AssignStakeReward ------------------------")
	// 0. setup
	err, proxy, impl := DeployAndUpgrade(pri1)
	if err != nil {
		panic("got error = " + err.Error())
	}
	fmt.Println("proxy = ", proxy)
	fmt.Println("impl = ", impl)
	p := NewProxy(api, proxy, impl)
	err2 := p.updateVerifier(pri1)
	if err2 != nil {
		panic("test assign stake reward failed: update verifier error: " + err2.Error())
	}

	_ = p.updateStakingParameter(pri1, "1000000000", "5000000000", "100000000000000000")


	ssn1 := "0x" + keytools.GetAddressFromPrivateKey(util.DecodeHex(pri1))
	ssn2 := "0x" + keytools.GetAddressFromPrivateKey(util.DecodeHex(pri2))

	if err0 := p.unpause(pri1); err0 != nil {
		panic("unpause with valid account failed")
	}
	// 1. assign ssn1, should fail
	err1 := p.assignStakeReward(pri1, ssn1, "50")
	if err1 == nil {
		panic("test assign stake with non-ssn reward failed")
	} else {
		res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
		sshmap := res["ssnlist"].(map[string]interface{})
		ssn := sshmap[ssn1]
		if ssn == nil {
			fmt.Println("test assign stake with non-ssn reward succeed")
		} else {
			panic("test assign stake with non-ssn reward failed")
		}
	}

	// 2. add ssn1, ssn2
	parameters := []contract2.Value{
		{
			VName: "ssnaddr",
			Type:  "ByStr20",
			Value: ssn1,
		},
		{
			VName: "stake_amount",
			Type:  "Uint128",
			Value: "100000000000",
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
			Value: "1000",
		},
	}
	args, _ := json.Marshal(parameters)
	if err3, output := ExecZli("contract", "call",
		"-k", pri1,
		"-a", proxy,
		"-t", "add_ssn_after_upgrade",
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
			fmt.Println("add ssn1 succeed")
		} else {
			panic("add ssn1 failed")
		}
	}

	parameters = []contract2.Value{
		{
			VName: "ssnaddr",
			Type:  "ByStr20",
			Value: ssn2,
		},
		{
			VName: "stake_amount",
			Type:  "Uint128",
			Value: "100000000000",
		},
		{
			VName: "rewards",
			Type:  "Uint128",
			Value: "0",
		},
		{
			VName: "urlraw",
			Type:  "String",
			Value: "devapiziiliqacom2",
		},
		{
			VName: "urlapi",
			Type:  "String",
			Value: "ziiliqacom2",
		},
		{
			VName: "buffered_deposit",
			Type:  "Uint128",
			Value: "0",
		},
	}
	args, _ = json.Marshal(parameters)
	if err3, output := ExecZli("contract", "call",
		"-k", pri1,
		"-a", proxy,
		"-t", "add_ssn_after_upgrade",
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
			fmt.Println("add ssn2 succeed")
		} else {
			panic("add ssn2 failed")
		}
	}

	// check totalstakedeposit
	res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
	totalstakedeposit := res["totalstakedeposit"].(string)
	if totalstakedeposit != "200000001000" {
		panic("check totalstakedeposit error")
	}

	// 3. reward ssn1
	err = p.assignStakeReward(pri1, ssn1, "50")
	if err != nil {
		panic("reward ssn1 failed: " + err.Error())
	}
	res = p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
	sshmap := res["ssnlist"].(map[string]interface{})
	ssn := sshmap[ssn1]
	if ssn == nil {
		panic("reward ssn1 failed")
	} else {
		arguments := ssn.(map[string]interface{})["arguments"].([]interface{})[2].(string)
		if arguments != "5000" {
			panic("reward ssn1 failed: state check error")
		} else {
			fmt.Println("reward ssn1 succeed")
		}
	}

	ssn = sshmap[ssn2]
	if ssn == nil {
		panic("_reward ssn2 failed")
	} else {
		arguments := ssn.(map[string]interface{})["arguments"].([]interface{})[2].(string)
		if arguments == "0" {
			fmt.Println("_reward ssn2 succeed")
		} else {
			panic("_reward ssn2 failed: state check error")
		}
	}

	// 4. reward ssn2
	err = p.assignStakeReward(pri1, ssn2, "50")
	if err != nil {
		panic("reward ssn2 failed: " + err.Error())
	}
	res = p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
	sshmap = res["ssnlist"].(map[string]interface{})
	ssn = sshmap[ssn2]
	if ssn == nil {
		panic("reward ssn2 failed")
	} else {
		arguments := ssn.(map[string]interface{})["arguments"].([]interface{})[2].(string)
		if arguments != "5000" {
			panic("reward ssn2 failed: state check error")
		} else {
			fmt.Println("reward ssn2 succeed")
		}
	}
}

func (p *Proxy) assignStakeReward(pri, ssn, percent string) error {
	proxy, _ := bech32.ToBech32Address(p.Addr)
	val := []Val{
		{
			Constructor: "SsnRewardShare",
			ArgTypes:    make([]interface{}, 0),
			Arguments:   []string{ssn, percent},
		},
	}
	parameters := []contract2.Value{
		{
			VName: "ssnreward_list",
			Type:  "List SsnRewardShare",
			Value: val,
		},
		{
			VName: "reward_blocknum",
			Type:  "Uint32",
			Value: "50000",
		},
	}

	args, _ := json.Marshal(parameters)

	if err2, output := ExecZli("contract", "call",
		"-k", pri,
		"-a", proxy,
		"-t", "assign_stake_reward",
		"-f", "true",
		"-r", string(args)); err2 != nil {
		return errors.New("call transition error: " + err2.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if success {
			return nil
		} else {
			return errors.New("transaction failed: " + tx)
		}
	}
}

type Val struct {
	Constructor string        `json:"constructor"`
	ArgTypes    []interface{} `json:"argtypes"`
	Arguments   []string      `json:"arguments"`
}
