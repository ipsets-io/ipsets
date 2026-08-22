package duckduckgo

import "github.com/ipsets-io/ipsets/provider"

const URL = "https://duckduckgo.com/duckduckbot.json"

func New() provider.Provider {
	return provider.PrefixList(provider.Meta{
		ID:        "duckduckgo",
		Name:      "DuckDuckGo",
		Homepage:  "https://duckduckgo.com/duckduckgo-help-pages/results/duckduckbot/",
		SourceURL: URL,
		Sets: []provider.Set{
			{ID: "duckduckbot", Name: "DuckDuckBot", Category: "crawler"},
		},
	})
}
