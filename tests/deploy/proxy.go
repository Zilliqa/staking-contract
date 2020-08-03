package deploy

import (
	"errors"
	"github.com/Zilliqa/gozilliqa-sdk/account"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"io/ioutil"
)

type Proxy struct {
	Code string
	Init []core.ContractValue
	Addr string
}

func NewProxy(key string) (*Proxy,error) {
	adminAddr := keytools.GetAddressFromPrivateKey(util.DecodeHex(key))
	code,_ := ioutil.ReadFile("./proxy.scilla")
	init := []core.ContractValue {
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "init_admin",
			Type:  "ByStr20",
			Value: "0x" + adminAddr,
		}, {
			VName: "init_implementation",
			Type:  "ByStr20",
			Value: "0x" + adminAddr,
		},
	}

	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)

	contract := contract2.Contract{
		Code:     string(code),
		Init:     init,
		Signer:   wallet,
	}

	tx,err := contract.DeployTo("isolated")
	if err != nil {
		return nil,err
	}
	tx.Confirm(tx.ID, 1000, 10, contract.Provider)
	if tx.Status == core.Confirmed {
		return &Proxy{
			Code: string(code),
			Init: init,
			Addr: tx.ContractAddress,
		},nil
	} else {
		return nil,errors.New("deploy failed")
	}
}
