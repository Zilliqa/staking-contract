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
	"strconv"
)

type Proxy struct {
	Code   string
	Init   []core.ContractValue
	Addr   string
	Bech32 string
	Wallet *account.Wallet
}

func (p *Proxy) WithdrawStakeAmount(ssn, amt string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssn",
			"ByStr20",
			ssn,
		},
		{
			"amt",
			"Uint128",
			amt,
		},
	}
	return p.Call("WithdrawStakeAmt", args, "0")
}

func (p *Proxy) UpdateComm(rate string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"new_rate",
			"Uint128",
			rate,
		},
	}
	return p.Call("UpdateComm", args, "0")
}

func (p *Proxy) AddDelegator(ssnaddr, deleg string, stakeAmount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			ssnaddr,
		},
		{
			"deleg",
			"ByStr20",
			deleg,
		},
		{
			"stake_amt",
			"Uint128",
			stakeAmount,
		},
	}
	return p.Call("UpdateDeleg", args, "0")
}

func (p *Proxy) PopulateTotalStakeAmt(amt string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"amt",
			"Uint128",
			amt,
		},
	}
	return p.Call("PopulateTotalStakeAmt", args, "0")
}

func (p *Proxy) DelegateStake(ssnaddr string, amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			ssnaddr,
		},
	}
	return p.Call("DelegateStake", args, amount)
}

func (p *Proxy) UpdateStakingParameters(min, delegmin string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"min_stake",
			"Uint128",
			min,
		},
		{
			"min_deleg_stake",
			"Uint128",
			delegmin,
		},
		{
			"max_comm_change_rate",
			"Uint128",
			"20",
		},
	}
	return p.Call("UpdateStakingParameters", args, "0")
}

func (p *Proxy) RemoveDelegator(ssnaddr, deleg string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			ssnaddr,
		},
		{
			"deleg",
			"ByStr20",
			deleg,
		},
	}
	return p.Call("Removedeleg", args, "0")
}

func (p *Proxy) RemoveSSN(addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			addr,
		},
	}

	return p.Call("RemoveSSN", args, "0")
}

func (p *Proxy) AddSSNAfterUpgrade(addr string, stakeAmt string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			addr,
		},
		{
			"stake_amt",
			"Uint128",
			stakeAmt,
		},
		{
			"rewards",
			"Uint128",
			"0",
		},
		{
			"name",
			"String",
			"fakename",
		},
		{
			"urlraw",
			"String",
			"fakeurl",
		},
		{
			"urlapi",
			"String",
			"fakeapi",
		},
		{
			"buff_deposit",
			"Uint128",
			"0",
		},
		{
			"comm",
			"Uint128",
			"0",
		},
		{
			"comm_rewards",
			"Uint128",
			"0",
		},
		{
			"rec_addr",
			"ByStr20",
			addr,
		},
	}

	return p.Call("AddSSNAfterUpgrade", args, "0")

}

func (p *Proxy) AddSSN(addr string, name string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			addr,
		},
		{
			"name",
			"String",
			name,
		},
		{
			"urlraw",
			"String",
			"fakeurl",
		},
		{
			"urlapi",
			"String",
			"fakeapi",
		},
		{
			"comm",
			"Uint128",
			"0",
		},
	}

	return p.Call("AddSSN", args, "0")
}

func (p *Proxy) UpdateReceiveAddr(newAddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{{
		"new_addr",
		"ByStr20",
		newAddr,
	}}
	return p.Call("UpdateReceivedAddr", args, "0")
}

type SSNRewardShare struct {
	SSNAddress       string
	RewardPercentage string
}

func (p *Proxy) AssignStakeReward(ssn, percent string) (*transaction.Transaction, error) {
	args := []core.ContractValue{{
		VName: "ssnreward_list",
		Type:  "List SsnRewardShare",
		Value: []core.ParamConstructor{
			{
				"SsnRewardShare",
				make([]interface{}, 0),
				[]string{ssn, percent},
			},
		},
	}, {
		"verifier_reward",
		"Uint128",
		"0",
	}}

	return p.Call("AssignStakeReward", args, "0")
}

func (p *Proxy) AssignStakeReward3(ssn1, percent1, ssn2, percent2 string) (*transaction.Transaction, error) {
	args := []core.ContractValue{{
		VName: "ssnreward_list",
		Type:  "List SsnRewardShare",
		Value: []core.ParamConstructor{
			{
				"SsnRewardShare",
				make([]interface{}, 0),
				[]string{ssn1, percent1},
			},
			{
				"SsnRewardShare",
				make([]interface{}, 0),
				[]string{ssn2, percent2},
			},
		},
	}, {
		"verifier_reward",
		"Uint128",
		"0",
	}}

	return p.Call("AssignStakeReward", args, "0")
}

func (p *Proxy) AssignStakeRewardBatch(ssn, percent string) []account.BatchSendingResult {
	args := []core.ContractValue{{
		VName: "ssnreward_list",
		Type:  "List SsnRewardShare",
		Value: []core.ParamConstructor{
			{
				"SsnRewardShare",
				make([]interface{}, 0),
				[]string{ssn, percent},
			},
		},
	}}
	return p.CallBatch("AssignStakeReward", args, "0")
}

func (p *Proxy) AssignStakeReward2(ssn1, percent1, ssn2, percent2 string) (*transaction.Transaction, error) {
	args := []core.ContractValue{{
		VName: "ssnreward_list",
		Type:  "List SsnRewardShare",
		Value: []core.ParamConstructor{
			{
				"SsnRewardShare",
				make([]interface{}, 0),
				[]string{ssn1, percent1},
			},
			{
				"SsnRewardShare",
				make([]interface{}, 0),
				[]string{ssn2, percent2},
			},
		},
	},
	}

	return p.Call("AssignStakeReward", args, "0")
}

func (p *Proxy) UpdateVerifier(addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{{
		"verif",
		"ByStr20",
		addr,
	}}
	return p.Call("UpdateVerifier", args, "0")

}

func (p *Proxy) WithdrawStakeRewards(addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{{
		"ssn_operator",
		"ByStr20",
		addr,
	}}
	return p.Call("WithdrawStakeRewards", args, "0")

}

func (p *Proxy) AddFunds(amount string) {
	args := []core.ContractValue{}
	_, err := p.Call("AddFunds", args, amount)
	if err != nil {
		log.Fatal("AddFunds failed")
	}
}

func (p *Proxy) WithdrawComm(ssnaddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{{
		"ssnaddr",
		"ByStr20",
		ssnaddr,
	}}
	return p.Call("WithdrawComm", args, "0")
}

func (p *Proxy) Unpause() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return p.Call("UnPause", args, "0")
}

func (p *Proxy) ClaimAdmin() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return p.Call("ClaimAdmin", args, "0")
}

func (p *Proxy) GetBalance() string {
	provider := provider2.NewProvider("https://zilliqa-isolated-server.zilliqa.com/")
	balAndNonce, _ := provider.GetBalance(p.Addr)
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

func (p *Proxy) CallBatch(transition string, params []core.ContractValue, amount string) []account.BatchSendingResult {
	var transactions []*transaction.Transaction

	for i := 0; i < 3; i++ {
		data := contract2.Data{
			Tag:    transition,
			Params: params,
		}
		txn := &transaction.Transaction{
			Version:      strconv.FormatInt(int64(util.Pack(1, 1)), 10),
			SenderPubKey: util.EncodeHex(p.Wallet.DefaultAccount.PublicKey),
			ToAddr:       p.Bech32,
			Amount:       "0",
			GasPrice:     "1000000000",
			GasLimit:     "40000",
			Code:         "",
			Data:         data,
			Priority:     false,
		}
		transactions = append(transactions, txn)
	}

	rpc := provider2.NewProvider("https://zilliqa-isolated-server.zilliqa.com/")
	p.Wallet.SignBatch(transactions, *rpc)
	return p.Wallet.SendBatch(transactions, *rpc)

}

func NewProxy(key string) (*Proxy, error) {
	adminAddr := keytools.GetAddressFromPrivateKey(util.DecodeHex(key))
	code, _ := ioutil.ReadFile("../contracts/proxy.scilla")
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
