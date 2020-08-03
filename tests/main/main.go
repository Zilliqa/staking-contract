package main

import (
	"Zilliqa/stake-test/deploy"
	"fmt"
)

func main() {
	// 1. deploy proxy
	proxy,err := deploy.NewProxy("e40afdc148c8f169613ba1bb2f9b15186cff6e1f5ad50ddc42aae7e5d51042bb")
	if err != nil {
		fmt.Println("deploy proxy error = ",err.Error())
		return
	}

	fmt.Println("deploy proxy succeed, address = ",proxy.Addr)

	// 2. deploy ssn list
	ssnlist,err1 := deploy.NewSSNList("e40afdc148c8f169613ba1bb2f9b15186cff6e1f5ad50ddc42aae7e5d51042bb",proxy.Addr)
	if err1 != nil {
		fmt.Println("deploy ssnlist error = ",err1.Error())
		return
	}

	fmt.Println("deploy proxy succeed, address = ",ssnlist.Addr)

}
