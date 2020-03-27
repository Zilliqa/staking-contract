package transitions

import (
	"encoding/json"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"io/ioutil"
	"os"
	"strings"
)

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
		"-f","true",
		"-t", "upgradeTo",
		"-r", string(args)); err2 != nil {
		return err2, "", ""
	} else {
		//fmt.Println(output)
		tx := strings.TrimSpace(strings.Split(output, "confirmed!")[1])
		fmt.Println("transaction id = ", tx)
		fmt.Println("------------------------ end upgrade ------------------------")
		return nil, proxyAddr, sshlistAddr
	}
}
