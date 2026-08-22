package apple

import "github.com/ipsets-io/ipsets/provider"

const URL = "https://search.developer.apple.com/applebot.json"

func New() provider.Provider {
	return provider.PrefixList(provider.Meta{
		ID:        "apple",
		Name:      "Apple",
		Homepage:  "https://support.apple.com/en-us/119829",
		SourceURL: URL,
		Sets: []provider.Set{
			{ID: "applebot", Name: "Applebot", Category: "crawler"},
		},
	})
}
