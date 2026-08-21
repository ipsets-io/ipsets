package perplexity

import "github.com/ipsets-io/ipsets/provider"

const URL = "https://www.perplexity.com/perplexitybot.json"

func New() provider.Provider {
	return provider.PrefixList(provider.Meta{
		ID:        "perplexity",
		Name:      "PerplexityBot",
		Homepage:  "https://docs.perplexity.ai/guides/bots",
		SourceURL: URL,
		Category:  "crawler",
	})
}
