package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func ExecZli(arg ...string) (error, string) {
	cmd := exec.Command("zli", arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err.Error()
	}
	defer stdout.Close()
	fmt.Println(cmd.String())
	if err := cmd.Start(); err != nil {
		return errors.New(fmt.Sprintf("exec error on command: \"zli %s\" %s", arg[0], err.Error())), ""
	}
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err.Error()
	}
	if err := cmd.Wait(); err != nil {
		return errors.New(fmt.Sprintf("exec error on wait: \"zli %s\" %s", arg[0], err.Error())), ""
	}
	return nil, string(opBytes)
}

// return proxy address and impl address
func DeployAndUpgrade(private string) (error, string, string) {
	// 1. deploy proxy contract with private key
	fmt.Println("------------------------ begin deploy proxy ------------------------")
	var proxyAddr string
	adminAddr := keytools.GetAddressFromPrivateKey(util.DecodeHex(private))
	proxyInit := []contract2.Value{
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
	data, _ := json.Marshal(proxyInit)
	jsonPath := "../contracts/proxy" + adminAddr + ".json"
	_ = ioutil.WriteFile(jsonPath, data, 0644)
	if err, output := ExecZli("contract", "deploy",
		"-k", private,
		"-c", "../contracts/proxy.scilla",
		"-i", jsonPath); err != nil {
		return err, "", ""
	} else {
		res := strings.Split(output, "contract address =  ")
		res = strings.Split(res[1], "{")
		res = strings.Split(res[0], "track")
		proxyAddr = strings.TrimSpace(res[0])
	}
	proxyBech32, _ := bech32.ToBech32Address(proxyAddr)
	fmt.Printf("proxy address = %s, bech32 = %s\n", proxyAddr, proxyBech32)
	fmt.Println("------------------------ end deploy proxy ------------------------")

	_ = os.Remove(jsonPath)

	// 2. deploy ssnlist contract with private key and proxy address
	fmt.Println("------------------------ begin deploy sshlist ------------------------")
	implInit := []contract2.Value{
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
			Value: "0x" + proxyAddr,
		},
	}
	data, _ = json.Marshal(implInit)
	fmt.Println(string(data))
	jsonPath = "../contracts/ssnlist" + adminAddr + ".json"
	_ = ioutil.WriteFile(jsonPath, data, 0644)
	var sshlistAddr string

	if err, output := ExecZli("contract", "deploy",
		"-k", private,
		"-c", "../contracts/ssnlist.scilla",
		"-i", jsonPath); err != nil {
		return err, "", ""
	} else {
		res := strings.Split(output, "contract address =  ")
		res = strings.Split(res[1], "{")
		res = strings.Split(res[0], "track")
		sshlistAddr = strings.TrimSpace(res[0])
	}
	sshlistBech32, _ := bech32.ToBech32Address(sshlistAddr)
	fmt.Printf("ssnlist address = %s, bech32 = %s\n", sshlistAddr, sshlistBech32)
	fmt.Println("------------------------ end deploy sshlist ------------------------")

	_ = os.Remove(jsonPath)

	// 3. upgrade to actual implement
	fmt.Println("------------------------ start upgrade ------------------------")
	parameters := []contract2.Value{
		{
			VName: "newImplementation",
			Type:  "ByStr20",
			Value: "0x" + sshlistAddr,
		},
	}
	args, _ := json.Marshal(parameters)
	fmt.Println(string(args))
	if err2, output := ExecZli("contract", "call",
		"-k", private,
		"-a", proxyBech32,
		"-t", "upgradeTo",
		"-r", string(args)); err2 != nil {
		return err2, "", ""
	} else {
		//fmt.Println(output)
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println(tx)
		fmt.Println("------------------------ end upgrade ------------------------")
		return nil, proxyAddr, sshlistAddr
	}

}

func main() {
	// 0. prepare two private keys, recommend to use your own, in case conflicts
	// we need the second private key, as we need to test non-admin permission or something similar
	pri1 := "55d256f225a0a552dc9c8158c87c460f92f9f18f4ae0f2ba104a69bf3ab7ed73"
	pri2 := "c25755f01577cb2d1c6a412ee8bfe2f98de0ed580844e5d7ae03bf0621c6b47e"
	// 1. make sure zli is already installed
	if err, output := ExecZli("-h"); err != nil {
		fmt.Println(err.Error())
		return
	} else {
		fmt.Println(output)
	}

	// setup default wallet (if there is no default one)
	if err, output := ExecZli("wallet", "init"); err != nil {
		// no need to panic here, only to make sure there is a wallet file
		// fmt.Println(err.Error())
	} else {
		fmt.Println(output)
	}
	if err, output := ExecZli("wallet", "echo"); err != nil {
		fmt.Println(err.Error())
		return
	} else {
		fmt.Println(output)
	}

	// deploy proxy and ssshlist, and do upgrade
	//  55d256f225a0a552dc9c8158c87c460f92f9f18f4ae0f2ba104a69bf3ab7ed73
	err, proxy, impl := DeployAndUpgrade(pri1)
	if err != nil {
		fmt.Println("got error = ", err.Error())
	}
	fmt.Println("proxy = ", proxy)
	fmt.Println("impl = ", impl)

	p := NewProxy("https://staking7-l2api.dev.z7a.xyz/", proxy, impl)
	//p := NewProxy("https://dev-api.zilliqa.com", "8adbd96a351cd89a286e69502ee2b641d2e03ac0", "b83c4bf4fc17054f53b86e9aa662610a2983f283")

	// test ChangeProxyAdmin
	p.ChangeProxyAdmin(pri1, pri2)

	// test UpdateAdmin
	p.UpdateAdmin(pri1, pri2)

	// test Pause
	p.PauseAndUnPause(pri1, pri2)

	// test UpdateVerifier
	p.UpdateVerifier(pri1, pri2)

	// test UpdateMinStake
	p.UpdateMinStake(pri1, pri2)

	// test UpdateMaxStake
	p.UpdateMaxStake(pri1, pri2)

	// test UpdateContractMaxStake
	p.UpdateContractMaxStake(pri1, pri2)
}

type Proxy struct {
	Addr        string
	ImplAddress string
	Provider    *provider.Provider
}

func NewProxy(url, proxy, impl string) *Proxy {
	p := provider.NewProvider(url)
	return &Proxy{
		Addr:        proxy,
		ImplAddress: impl,
		Provider:    p,
	}
}

func (p *Proxy) UpdateContractMaxStake(valid, invalid string) {
	err := p.updateContractMaxStake(invalid)
	if err == nil {
		panic("update contract max stake with invalid key failed")
	}

	fmt.Println("update contract max stake with invalid key succeed")

	err2 := p.updateContractMaxStake(valid)
	if err2 != nil {
		panic("update contract max stake with valid key failed")
	}

	fmt.Println("update contract max stake with valid key succeed")

}

func (p *Proxy) updateContractMaxStake(private string) error {
	proxy, _ := bech32.ToBech32Address(p.Addr)
	stakeNum := fmt.Sprintf("%d", rand.Int())
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
		"-t", "update_contractmaxstake",
		"-r", string(args)); err2 != nil {
		return errors.New("call transition error: " + err2.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if success {
			res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
			minstake := res["contractmaxstake"].(string)
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
	stakeNum := fmt.Sprintf("%d", rand.Int())
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

func (p *Proxy) UpdateMinStake(valid, invalid string) {
	err := p.updateMinStake(invalid)
	if err == nil {
		panic("update min stake with invalid key failed")
	}

	fmt.Println("update min stake with invalid key succeed")

	err2 := p.updateMinStake(valid)
	if err2 != nil {
		panic("update min stake with valid key failed")
	}

	fmt.Println("update min stake with valid key succeed")

}

func (p *Proxy) updateMinStake(private string) error {
	proxy, _ := bech32.ToBech32Address(p.Addr)
	stakeNum := fmt.Sprintf("%d", rand.Int())
	parameters := []contract2.Value{
		{
			VName: "min_stake",
			Type:  "Uint128",
			Value: stakeNum,
		},
	}
	args, _ := json.Marshal(parameters)
	if err2, output := ExecZli("contract", "call",
		"-k", private,
		"-a", proxy,
		"-t", "update_minstake",
		"-r", string(args)); err2 != nil {
		return errors.New("call transition error: " + err2.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		payload := p.Provider.GetTransaction(tx).Result.(map[string]interface{})
		receipt := payload["receipt"].(map[string]interface{})
		success := receipt["success"].(bool)
		if success {
			res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
			minstake := res["minstake"].(string)
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

func (p *Proxy) updateVerifier(private string) error {
	proxy, _ := bech32.ToBech32Address(p.Addr)
	operator := "0x" + keytools.GetAddressFromPrivateKey(util.DecodeHex(private))
	parameters := []contract2.Value{
		{
			VName: "verif",
			Type:  "ByStr20",
			Value: operator,
		},
	}
	args, _ := json.Marshal(parameters)
	if err2, output := ExecZli("contract", "call",
		"-k", private,
		"-a", proxy,
		"-t", "update_verifier",
		"-r", string(args)); err2 != nil {
		return errors.New("call transition error: " + err2.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
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
			if arg == operator {
				return nil
			} else {
				return errors.New("update state failed")
			}
		} else {
			return errors.New("transaction failed")
		}
	}
}

func (p *Proxy) PauseAndUnPause(valid, invalid string) {
	fmt.Println("------------------------ begin pause unpause ------------------------")
	if err2 := p.pause(invalid); err2 == nil {
		panic("pause with invalid account failed")
	} else {
		if err2.Error() == "failed" {
			fmt.Println("pause with invalid account succeed")
		} else {
			panic("pause with invalid account failed: " + err2.Error())
		}
	}

	err := p.pause(valid)
	if err != nil {
		panic("pause with valid account failed: " + err.Error())
	}
	fmt.Println("pause with valid account succeed")

	if err3 := p.unpause(invalid); err3 == nil {
		panic("unpause with invalid account failed")
	} else {
		if err3.Error() == "failed" {
			fmt.Println("unpause with invalid account succeed")
		} else {
			panic("unpause with invalid account failed: " + err3.Error())
		}
	}
	err4 := p.unpause(valid)
	if err4 != nil {
		panic("unpause with valid account failed: " + err4.Error())
	}
	fmt.Println("unpause with valid account succeed")
	fmt.Println("------------------------ end pause unpause ------------------------")

}

func (p *Proxy) pause(private string) error {
	proxy, _ := bech32.ToBech32Address(p.Addr)

	if err2, _ := ExecZli("contract", "call",
		"-k", private,
		"-a", proxy,
		"-t", "pause",
		"-r", "[]"); err2 != nil {
		return errors.New("call transition error: " + err2.Error())
	} else {
		res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
		paused := res["paused"].(map[string]interface{})
		constructor := paused["constructor"].(string)
		if constructor == "True" {
			return nil
		} else {
			return errors.New("failed")
		}
	}
}

func (p *Proxy) unpause(private string) error {
	proxy, _ := bech32.ToBech32Address(p.Addr)

	if err2, _ := ExecZli("contract", "call",
		"-k", private,
		"-a", proxy,
		"-t", "unpause",
		"-r", "[]"); err2 != nil {
		return errors.New("call transition error: " + err2.Error())
	} else {
		res := p.Provider.GetSmartContractState(p.ImplAddress).Result.(map[string]interface{})
		paused := res["paused"].(map[string]interface{})
		constructor := paused["constructor"].(string)
		if constructor == "False" {
			return nil
		} else {
			return errors.New("failed")
		}
	}
}

func (p *Proxy) UpdateAdmin(oldPrivateKey, newPrivateKey string) {
	fmt.Println("------------------------ begin update admin ------------------------")
	err := p.updateAdmin(oldPrivateKey, newPrivateKey)
	if err != nil {
		panic("update admin with admin permission succeed" + err.Error())
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
		"-r", string(args)); err2 != nil {
		return errors.New("call transition error: " + err2.Error())
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
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
		"-r", string(args)); err2 != nil {
		return err2
	} else {
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
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

func (p *Proxy) extraEventName(events []interface{}) string {
	logs := events[0]
	log := logs.(map[string]interface{})
	event := log["_eventname"].(string)
	return event
}
