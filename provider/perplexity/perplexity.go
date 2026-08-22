package perplexity

import "github.com/ipsets-io/ipsets/provider"

const URL = "https://www.perplexity.com/perplexitybot.json"

func New() provider.Provider {
	return provider.PrefixList(provider.Meta{
		ID:        "perplexity",
		Name:      "Perplexity",
		Homepage:  "https://docs.perplexity.ai/guides/bots",
		SourceURL: URL,
		Sets: []provider.Set{
			{ID: "perplexitybot", Name: "PerplexityBot", Category: "crawler"},
		},
	})
}
