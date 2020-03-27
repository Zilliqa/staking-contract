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

func (p *Proxy) UpdateAdmin(oldPrivateKey, newPrivateKey string) {
	fmt.Println("------------------------ begin update admin ------------------------")
	err := p.updateAdmin(oldPrivateKey, newPrivateKey)
	if err != nil {
		panic("update admin with admin permission failed" + err.Error())
	}
	fmt.Println("update admin with admin permission succeed")
	err2 := p.updateAdmin(oldPrivateKey, newPrivateKey)
	if err2 == nil {
		panic("update admin without admin permission failed")
	}
	fmt.Println("update admin without admin permission succeed")
	err3 := p.updateAdmin(newPrivateKey, oldPrivateKey)
	if err3 != nil {
		panic("revert admin failed")
	}
	fmt.Println("revert admin succeed")
	fmt.Println("------------------------ end update admin ------------------------")
}

func (p *Proxy) updateAdmin(oldPrivateKey, newPrivateKey string) error {
	newAddr := "0x" + keytools.GetAddressFromPrivateKey(util.DecodeHex(newPrivateKey))
	proxy, _ := bech32.ToBech32Address(p.Addr)
	parameters := []contract2.Value{
		{
			VName: "admin",
			Type:  "ByStr20",
			Value: newAddr,
		},
	}

	args, _ := json.Marshal(parameters)
	if err2, output := ExecZli("contract", "call",
		"-k", oldPrivateKey,
		"-a", proxy,
		"-t", "update_admin",
		"-f","true",
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
			contractAdmin := res["contractadmin"].(string)
			if contractAdmin == newAddr {
				return nil
			} else {
				return errors.New("update state error")
			}
		} else {
			return errors.New("transaction failed")
		}
	}
}
