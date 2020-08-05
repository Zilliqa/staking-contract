package deploy

import (
	"encoding/json"
	"errors"
	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"io/ioutil"
	"log"
)

type Proxy struct {
	Code   string
	Init   []core.ContractValue
	Addr   string
	Bech32 string
	Wallet *account.Wallet
}

func (p *Proxy) UpdateVerifier(addr string) (*transaction.Transaction, error){
	args := []core.ContractValue{{
		"verif",
		"ByStr20",
		addr,
	}}
	return p.Call("UpdateVerifier",args,"0")

}

func (p *Proxy) AddFunds(amount string) {
	args := []core.ContractValue{}
	_, err := p.Call("AddFunds", args,amount)
	if err != nil {
		log.Fatal("AddFunds failed")
	}
}

func (p *Proxy) Unpause() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return p.Call("UnPause", args,"0")
}

func (p *Proxy) GetBalance() string {
	provider := provider2.NewProvider("https://zilliqa-isolated-server.zilliqa.com/")
	balAndNonce,_ := provider.GetBalance(p.Addr)
	return balAndNonce.Balance
}

func (p *Proxy) UpdateWallet(newKey string) {
	wallet := account.NewWallet()
	wallet.AddByPrivateKey(newKey)
	p.Wallet = wallet
}

func (p *Proxy) LogContractStateJson() {
	provider := provider2.NewProvider("https://zilliqa-isolated-server.zilliqa.com/")
	rsp, _ := provider.GetSmartContractState(p.Addr)
	j, _ := json.Marshal(rsp)
	log.Println(string(j))
}

func (p *Proxy) Call(transition string, params []core.ContractValue, amount string) (*transaction.Transaction, error) {
	contract := contract2.Contract{
		Address: p.Bech32,
		Signer:  p.Wallet,
	}

	tx, err := contract.CallFor(transition, params, false, amount, "isolated")
	if err != nil {
		return tx, err
	}
	tx.Confirm(tx.ID, 1000, 3, contract.Provider)
	if tx.Status != core.Confirmed {
		return tx, errors.New("transaction didn't get confirmed")
	}
	if !tx.Receipt.Success {
		return tx, errors.New("transaction failed")
	}
	return tx, nil
}

func NewProxy(key string) (*Proxy, error) {
	adminAddr := keytools.GetAddressFromPrivateKey(util.DecodeHex(key))
	code, _ := ioutil.ReadFile("./proxy.scilla")
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
			VName: "init_implementation",
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

	tx, err := contract.DeployTo("isolated")
	if err != nil {
		return nil, err
	}
	tx.Confirm(tx.ID, 1000, 10, contract.Provider)

	b32, _ := bech32.ToBech32Address(tx.ContractAddress)
	if tx.Status == core.Confirmed {
		return &Proxy{
			Code:   string(code),
			Init:   init,
			Addr:   tx.ContractAddress,
			Wallet: wallet,
			Bech32: b32,
		}, nil
	} else {
		return nil, errors.New("deploy failed")
	}
}
