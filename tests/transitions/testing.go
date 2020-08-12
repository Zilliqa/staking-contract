package transitions

import (
	"encoding/json"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"log"
	"strings"
)

const key1 = "e40afdc148c8f169613ba1bb2f9b15186cff6e1f5ad50ddc42aae7e5d51042bb"
const addr1 = "29cf16563fac1ad1596dfe6f333978fece9706ec"
const key2 = "8732034b0c895564d966e3df6968205211c7a2f0140b77c9e13de10c1ce77873"
const addr2 = "e2cd74983c7a3487af3a133a3bf4e7dd76f5d928"
const key3 = "70c57a0a1f9a0e2c9192f28279a491bcb30a7d0ada87eab9aa0b6afad3f31c91"
const addr3 = "8bdc7e9064f3963654967fa28976aac98f002a58"
const key4 = "243d302e971f7469cb20cc4d37c4629f0c22860667370b4d1130ae4ab1a5f4f9"
const addr4 = "6e081b8cca40c585d6d69f9643faf1a545d13d63"

type Testing struct {
}

func NewTesting() *Testing {
	return &Testing{}
}

func (t *Testing) LogStart(tag string) {
	log.Printf("start to test %s\n", tag)
}

func (t *Testing) LogEnd(tag string) {
	log.Printf("end to test %s\n", tag)
}

func (t *Testing) LogError(tag string, err error) {
	log.Fatalf("failed at %s, err = %s\n", tag, err.Error())
}

func (t *Testing) AssertContain(s1, s2 string) {
	if !strings.Contains(s1, s2) {
		log.Println(s1)
		log.Println(s2)
		log.Fatal("assert failed")
	}
}

func (t *Testing) AssertError(err error) {
	if err == nil {
		log.Fatal("assert error failed")
	}
}

func (t *Testing) GetReceiptString(tnx *transaction.Transaction) string {
	receipt, _ := json.Marshal(tnx.Receipt)
	return string(receipt)
}
