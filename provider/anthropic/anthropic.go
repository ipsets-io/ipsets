package anthropic

import "github.com/ipsets-io/ipsets/provider"

const URL = "https://claude.com/crawling/bots.json"

func New() provider.Provider {
	return provider.PrefixList(provider.Meta{
		ID:        "anthropic",
		Name:      "Anthropic",
		SourceURL: URL,
		Sets: []provider.Set{
			{ID: "bots", Name: "ClaudeBot, Claude-User, Claude-SearchBot", Category: "crawler"},
		},
	})
}
