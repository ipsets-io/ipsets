package ipsets

import (
	"github.com/ipsets-io/ipsets/provider"
	"github.com/ipsets-io/ipsets/provider/anthropic"
	"github.com/ipsets-io/ipsets/provider/apple"
	"github.com/ipsets-io/ipsets/provider/aws"
	"github.com/ipsets-io/ipsets/provider/azure"
	"github.com/ipsets-io/ipsets/provider/bing"
	"github.com/ipsets-io/ipsets/provider/bunny"
	"github.com/ipsets-io/ipsets/provider/cloudflare"
	"github.com/ipsets-io/ipsets/provider/datadog"
	"github.com/ipsets-io/ipsets/provider/duckduckgo"
	"github.com/ipsets-io/ipsets/provider/fastly"
	"github.com/ipsets-io/ipsets/provider/github"
	"github.com/ipsets-io/ipsets/provider/google"
	"github.com/ipsets-io/ipsets/provider/openai"
	"github.com/ipsets-io/ipsets/provider/oracle"
	"github.com/ipsets-io/ipsets/provider/perplexity"
	"github.com/ipsets-io/ipsets/provider/pingdom"
	"github.com/ipsets-io/ipsets/provider/sentry"
	"github.com/ipsets-io/ipsets/provider/tor"
)

func Providers() []provider.Provider {
	return []provider.Provider{
		anthropic.New(),
		apple.New(),
		aws.New(),
		azure.New(),
		bing.New(),
		bunny.New(),
		cloudflare.New(),
		datadog.New(),
		duckduckgo.New(),
		fastly.New(),
		github.New(),
		google.New(),
		openai.New(),
		oracle.New(),
		perplexity.New(),
		pingdom.New(),
		sentry.New(),
		tor.New(),
	}
}
