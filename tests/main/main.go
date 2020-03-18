package main

import (
	"Zilliqa/stake-test/transitions"
	"fmt"
)

func main() {
	// 0. prepare two private keys, recommend to use your own, in case conflicts
	// we need the second private key, as we need to test non-admin permission or something similar
	pri1 := "55d256f225a0a552dc9c8158c87c460f92f9f18f4ae0f2ba104a69bf3ab7ed73"
	pri2 := "c25755f01577cb2d1c6a412ee8bfe2f98de0ed580844e5d7ae03bf0621c6b47e"
	api := "https://staking7-l2api.dev.z7a.xyz/"
	// 1. make sure zli is already installed
	if err, output := transitions.ExecZli("-h"); err != nil {
		fmt.Println(err.Error())
		return
	} else {
		fmt.Println(output)
	}

	// setup default wallet (if there is no default one)
	if err, output := transitions.ExecZli("wallet", "init"); err != nil {
		// no need to panic here, only to make sure there is a wallet file
		// fmt.Println(err.Error())
	} else {
		fmt.Println(output)
	}
	if err, output := transitions.ExecZli("wallet", "echo"); err != nil {
		fmt.Println(err.Error())
		return
	} else {
		fmt.Println(output)
	}

	err, proxy, impl := transitions.DeployAndUpgrade(pri1)
	if err != nil {
		panic("got error = " + err.Error())
	}
	fmt.Println("proxy = ", proxy)
	fmt.Println("impl = ", impl)

	p := transitions.NewProxy(api, proxy, impl)

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

	// test deposit
	p.TransferFunds(pri1)

	// test AddSSN
	transitions.TestAddSSN(pri1, pri2, api)
}
