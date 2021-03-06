package deploy

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/util"
)

type SSNList struct {
	Code string
	Init []core.ContractValue
	Addr string
}

func (s *SSNList) LogContractStateJson() string {
	provider := provider2.NewProvider("https://zilliqa-isolated-server.zilliqa.com/")
	rsp, _ := provider.GetSmartContractState(s.Addr)
	j, _ := json.Marshal(rsp)
	s.LogPrettyStateJson(rsp)
	return string(j)
}

func (s *SSNList) LogPrettyStateJson(data interface{}) {
	j, _ := json.MarshalIndent(data, "", "   ")
	log.Println(string(j))
}

func (s *SSNList) GetBalance() string {
	provider := provider2.NewProvider("https://zilliqa-isolated-server.zilliqa.com/")
	balAndNonce, _ := provider.GetBalance(s.Addr)
	return balAndNonce.Balance
}

func NewSSNList(key string, proxy string) (*SSNList, error) {
	code, _ := ioutil.ReadFile("./ssnlist.scilla")
	adminAddr := keytools.GetAddressFromPrivateKey(util.DecodeHex(key))

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "init_admin",
			Type:  "ByStr20",
			Value: "0x" + adminAddr,
		}, {
			VName: "init_proxy_address",
			Type:  "ByStr20",
			Value: "0x" + proxy,
		},
		{
			VName: "init_gzil_address",
			Type:  "ByStr20",
			Value: "0x" + adminAddr,
		},
	}

	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)

	contract := contract2.Contract{
		Code:   string(code),
		Init:   init,
		Signer: wallet,
	}

	tx, err := contract.DeployTo("isolated", "50000")
	if err != nil {
		return nil, err
	}
	tx.Confirm(tx.ID, 1000, 10, contract.Provider)
	if tx.Status == core.Confirmed {
		return &SSNList{
			Code: string(code),
			Init: init,
			Addr: tx.ContractAddress,
		}, nil
	} else {
		return nil, errors.New("deploy failed")
	}
}
