package transitions

import (
	"Zilliqa/stake-test/deploy"
	"log"
)

func (t *Testing) LoadRemoteProxySSN() (*deploy.Proxy, *deploy.SSNList) {

	proxyAddress := "checksum_address_without_0x"
	implAddress := "checksum_address_without_0x"

	log.Println("start to load remote proxy")
	log.Println("proxy: " + proxyAddress)
	proxy, _ := deploy.LoadRemoteProxy(key1, proxyAddress)

	impl, _ := deploy.LoadRemoteSSN(key1, proxyAddress, implAddress)

	return proxy, impl
}
