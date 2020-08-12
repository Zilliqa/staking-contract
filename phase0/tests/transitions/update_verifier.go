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

func (p *Proxy) UpdateVerifier(valid, invalid string) {
	err := p.updateVerifier(valid)
	if err != nil {
		panic("update verifier with admin failed: " + err.Error())
	}
	fmt.Println("update verifier with admin succeed")

	err2 := p.updateVerifier(invalid)
	if err2 == nil {
		panic("update verifier with invalid admin failed")
	}
	fmt.Println("update verifier with invalid admin succeed")
}

func (p *Proxy) updateVerifierTo(private string, newVerifier string) error {
	proxy, _ := bech32.ToBech32Address(p.Addr)

	parameters := []contract2.Value{
		{
			VName: "verif",
			Type:  "ByStr20",
			Value: newVerifier,
		},
	}
	args, _ := json.Marshal(parameters)
	if err2, output := ExecZli("contract", "call",
		"-k", private,
		"-a", proxy,
		"-t", "update_verifier",
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
			res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
			verifier := res["verifier"].(map[string]interface{})
			arguments := verifier["arguments"].([]interface{})[0]
			fmt.Println(arguments)
			if arguments == nil {
				return errors.New("verifier is none")
			}
			arg := arguments.(string)
			if arg == newVerifier {
				return nil
			} else {
				return errors.New("update state failed")
			}
		} else {
			return errors.New("transaction failed")
		}
	}
}

func (p *Proxy) updateVerifier(private string) error {
	operator := "0x" + keytools.GetAddressFromPrivateKey(util.DecodeHex(private))
	return p.updateVerifierTo(private, operator)
}
