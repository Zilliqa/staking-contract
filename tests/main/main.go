package main

import (
	"Zilliqa/stake-test/transitions"
	"fmt"
	"os"
	"time"
)

func main() {
	// 0. prepare two private keys, recommend to use your own, in case conflicts
	// we need the second private key, as we need to test non-admin permission or something similar
	fromTime := time.Now()
	args := os.Args[1:]
	if len(args) != 3 {
		panic("parameters error")
	}
	pri1 := args[0]
	pri2 := args[1]
	api := args[2]

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

	p.UpdateStakingParameter(pri1,pri2)

	// test deposit
	p.TransferFundsAndDrainBalance(pri1, pri2, "5000")

	// test AddSSN
	transitions.TestAddSSN(pri1, pri2, api)

	transitions.TestStakeDeposit(pri1, pri2, api)

	transitions.TestWithdrawStakeRewards(pri1, pri2, api)

	transitions.TestAssignStakeReward(pri1, pri2, api)

	transitions.TestWithdrawAmount(pri1, pri2, api)
	
	transitions.TestRemoveSSN(pri1,pri2,api)

	endTime := time.Now()
	interval := endTime.Sub(fromTime).Minutes()
	fmt.Printf("The whole test cost %f minutes, oh my!\n", interval)
}
