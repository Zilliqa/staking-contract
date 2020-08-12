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

func (p *Proxy) ChangeProxyAdmin(oldPrivateKey, newPrivateKey string) {
	// change admin to fnewPrivateKey
	fmt.Println("------------------------ start ChangeProxyAdmin ------------------------")
	// test one have admin permission
	err2 := p.changeProxyAdmin(oldPrivateKey, newPrivateKey)
	if err2 != nil {
		panic("ChangeProxyAdmin failed: " + err2.Error())
	} else {
		fmt.Println("ChangeProxyAdmin test one have admin permission succeed")
	}

	// test one does not have admin permission
	err3 := p.changeProxyAdmin(oldPrivateKey, newPrivateKey)

	if err3 == nil {
		panic("ChangeProxyAdmin failed: ")
	} else {
		fmt.Println("ChangeProxyAdmin test one does not have admin permission succeed")
	}

	// revert admin permission
	err4 := p.changeProxyAdmin(newPrivateKey, oldPrivateKey)
	if err4 != nil {
		panic("ChangeProxyAdmin failed: " + err4.Error())
	} else {
		fmt.Println("ChangeProxyAdmin revert admin account succeed")
	}
	fmt.Println("------------------------ end ChangeProxyAdmin ------------------------")

}

func (p *Proxy) changeProxyAdmin(oldPrivateKey, newPrivateKey string) error {
	newAddr := "0x" + keytools.GetAddressFromPrivateKey(util.DecodeHex(newPrivateKey))
	proxy, _ := bech32.ToBech32Address(p.Addr)
	parameters := []contract2.Value{
		{
			VName: "newAdmin",
			Type:  "ByStr20",
			Value: newAddr,
		},
	}
	args, _ := json.Marshal(parameters)
	if err2, output := ExecZli("contract", "call",
		"-k", oldPrivateKey,
		"-a", proxy,
		"-t", "changeProxyAdmin",
		"-f","true",
		"-r", string(args)); err2 != nil {
		return err2
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		events := receipt["event_logs"]
		if success {
			event := p.extraEventName(events.([]interface{}))
			if event == "changeAdmin FailedNotAdmin" {
				return errors.New("ChangeProxyAdmin failed")
			}
			res := p.Provider.GetSmartContractState(p.Addr).Result.(map[string]interface{})
			admin := res["admin"].(string)
			if newAddr == admin {
				return nil
			} else {
				return errors.New("ChangeProxyAdmin failed, state not equal")
			}

		} else {
			return errors.New("transaction failed")
		}
	}
}
