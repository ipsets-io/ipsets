package ipsets

import (
	"github.com/ipsets-io/ipsets/provider"
	"github.com/ipsets-io/ipsets/provider/anthropic"
	"github.com/ipsets-io/ipsets/provider/apple"
	"github.com/ipsets-io/ipsets/provider/aws"
	"github.com/ipsets-io/ipsets/provider/bing"
	"github.com/ipsets-io/ipsets/provider/cloudflare"
	"github.com/ipsets-io/ipsets/provider/duckduckgo"
	"github.com/ipsets-io/ipsets/provider/fastly"
	"github.com/ipsets-io/ipsets/provider/github"
	"github.com/ipsets-io/ipsets/provider/google"
	"github.com/ipsets-io/ipsets/provider/openai"
	"github.com/ipsets-io/ipsets/provider/perplexity"
)

func Providers() []provider.Provider {
	return []provider.Provider{
		anthropic.New(),
		apple.New(),
		aws.New(),
		bing.New(),
		cloudflare.New(),
		duckduckgo.New(),
		fastly.New(),
		github.New(),
		google.New(),
		openai.New(),
		perplexity.New(),
	}
}
