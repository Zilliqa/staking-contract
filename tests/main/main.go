package main

import "Zilliqa/stake-test/transitions"

func main() {
	t := transitions.NewTesting()
	t.UpgradeTo()
}
