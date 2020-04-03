package transitions

import (
	"errors"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
)

func (p *Proxy) PauseAndUnPause(valid, invalid string) {
	fmt.Println("------------------------ begin pause unpause ------------------------")
	if err1 := p.unpause(valid); err1 != nil {
		panic("unpause with valid account failed")
	}
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
		"-f","true",
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
		"-f","true",
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
