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

type SSNList struct {
	Code string
	Init []core.ContractValue
	Addr string
}

func NewSSNList(key string, proxy string) (*SSNList,error) {
	code,_ := ioutil.ReadFile("./ssnlist.scilla")
	adminAddr := keytools.GetAddressFromPrivateKey(util.DecodeHex(key))

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
			VName: "proxy_address",
			Type:  "ByStr20",
			Value: "0x" + proxy,
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
		return &SSNList{
			Code: string(code),
			Init: init,
			Addr: tx.ContractAddress,
		},nil
	} else {
		return nil,errors.New("deploy failed")
	}
}
