package bing

import "github.com/ipsets-io/ipsets/provider"

const URL = "https://www.bing.com/toolbox/bingbot.json"

func New() provider.Provider {
	return provider.PrefixList(provider.Meta{
		ID:        "bing",
		Name:      "Bingbot",
		Homepage:  "https://www.bing.com/webmasters/help/how-to-verify-bingbot-3905dc26",
		SourceURL: URL,
		Category:  "crawler",
	})
}
