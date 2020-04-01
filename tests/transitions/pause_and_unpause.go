package transitions

import (
	"errors"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
)

func (p *Proxy) PauseAndUnPause(valid, invalid string) {
	fmt.Println("------------------------ begin pause unpause ------------------------")
	err := p.unpause(valid)
	if err != nil {
		panic("unpause with valid account failed: " + err.Error())
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

	err3 := p.pause(valid)
	if err3 != nil {
		panic("pause with valid account failed: " + err3.Error())
	}
	fmt.Println("pause with valid account succeed")

	if err4 := p.unpause(invalid); err4 == nil {
		panic("unpause with invalid account failed")
	} else {
		if err4.Error() == "failed" {
			fmt.Println("unpause with invalid account succeed")
		} else {
			panic("unpause with invalid account failed: " + err4.Error())
		}
	}
	err5 := p.unpause(valid)
	if err5 != nil {
		panic("unpause with valid account failed: " + err5.Error())
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
