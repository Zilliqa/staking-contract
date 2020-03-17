package transitions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"strings"
)

func (p *Proxy) UpdateMaxStake(valid, invalid string) {
	err := p.updateMaxStake(invalid)
	if err == nil {
		panic("update max stake with invalid key failed")
	}

	fmt.Println("update max stake with invalid key succeed")

	err2 := p.updateMaxStake(valid)
	if err2 != nil {
		panic("update max stake with valid key failed")
	}

	fmt.Println("update max stake with valid key succeed")

}

func (p *Proxy) updateMaxStake(private string) error {
	proxy, _ := bech32.ToBech32Address(p.Addr)
	//stakeNum := fmt.Sprintf("%d", rand.Int())
	stakeNum := "5000000000"
	parameters := []contract2.Value{
		{
			VName: "max_stake",
			Type:  "Uint128",
			Value: stakeNum,
		},
	}
	args, _ := json.Marshal(parameters)
	if err2, output := ExecZli("contract", "call",
		"-k", private,
		"-a", proxy,
		"-t", "update_maxstake",
		"-f", "true",
		"-r", string(args)); err2 != nil {
		return errors.New("call transition error: " + err2.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if success {
			res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
			minstake := res["maxstake"].(string)
			if minstake == stakeNum {
				return nil
			} else {
				return errors.New("state failed")
			}

		} else {
			return errors.New("transaction failed")
		}
	}
}
