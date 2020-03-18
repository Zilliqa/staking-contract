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
func TestAssignStakeReward(pri1, pri2,api string) {
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

	ssnaddr := "0x" + keytools.GetAddressFromPrivateKey(util.DecodeHex(pri2))

	// 1. assign ssn, should fail
	err1 := p.assignStakeReward(pri1, ssnaddr)
	if err1 == nil {
		panic("test assign stake with non-ssn reward failed")
	} else {
		fmt.Println("test assign stake with non-ssn reward succeed")
	}

}

func (p *Proxy) assignStakeReward(pri, ssn string) error {
	proxy, _ := bech32.ToBech32Address(p.Addr)
	val := []Val{
		{
			Constructor: "SsnRewardShare",
			ArgTypes:    make([]interface{}, 0),
			Arguments:   []string{ssn},
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
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if success {
			return nil
		} else {
			return errors.New("transaction failed")
		}
	}
}

type Val struct {
	Constructor string        `json:"constructor"`
	ArgTypes    []interface{} `json:"argtypes"`
	Arguments   []string      `json:"arguments"`
}
