package transitions

import (
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"log"
)

func (t *Testing) UpdateVerifier() {
	t.LogStart("UpdateVerifier")

	proxy,ssnlist := t.DeployAndUpgrade()

	// as admin, update verifier
	fakeVerifier := "0x82b142aa3d6733f8373477eca5cb2f35f240928f"
	ssnlist.LogContractStateJson()
	args := []core.ContractValue{{
		"verif",
		"ByStr20",
		fakeVerifier,
	}}
	tnx, err := proxy.Call("UpdateVerifier", args)
	if err != nil {
		t.LogError("UpdateVerifier",err)
	}
	receipt := t.GetReceiptString(tnx)
	log.Println(receipt)
	t.AssertContain(receipt, fakeVerifier)
	ssnlist.LogContractStateJson()


	// as non admin, update verifier
	fakeVerifier = "0xc61556c0762bd6ffd05258e083fdf70aa7537c3b"
	args = []core.ContractValue{{
		"verif",
		"ByStr20",
		fakeVerifier,
	}}
	proxy.UpdateWallet(key2)
	tnx, err1 := proxy.Call("UpdateVerifier", args)
	t.AssertError(err1)
	receipt = t.GetReceiptString(tnx)
	t.AssertContain(receipt,"Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -3))])")
	log.Println(receipt)

	t.LogEnd("UpdateVerifier")
}