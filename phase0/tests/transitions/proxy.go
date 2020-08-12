package transitions

import "github.com/Zilliqa/gozilliqa-sdk/provider"

type Proxy struct {
	Addr        string
	ImplAddress string
	Provider    *provider.Provider
}

func NewProxy(url, proxy, impl string) *Proxy {
	p := provider.NewProvider(url)
	return &Proxy{
		Addr:        proxy,
		ImplAddress: impl,
		Provider:    p,
	}
}

func (p *Proxy) extraEventName(events []interface{}) string {
	logs := events[0]
	log := logs.(map[string]interface{})
	event := log["_eventname"].(string)
	return event
}
