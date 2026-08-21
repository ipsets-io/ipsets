package google

import (
	"context"
	"net/http"

	"github.com/ipsets-io/ipsets/provider"
)

const (
	URLGoog     = "https://www.gstatic.com/ipranges/goog.json"
	URLCloud    = "https://www.gstatic.com/ipranges/cloud.json"
	crawlerBase = "https://developers.google.com/static/search/apis/ipranges/"
)

var sets = []provider.Set{
	{ID: "goog", Name: "Google (all ranges)", Category: "cloud", Source: URLGoog},
	{ID: "cloud", Name: "Google Cloud", Category: "cloud", Source: URLCloud},
	{ID: "googlebot", Name: "Googlebot", Category: "crawler", Source: crawlerBase + "googlebot.json"},
	{ID: "special-crawlers", Name: "Google special crawlers", Category: "crawler", Source: crawlerBase + "special-crawlers.json"},
	{ID: "user-triggered-fetchers", Name: "Google user-triggered fetchers", Category: "crawler", Source: crawlerBase + "user-triggered-fetchers.json"},
	{ID: "user-triggered-fetchers-google", Name: "Google user-triggered fetchers (Google IPs)", Category: "crawler", Source: crawlerBase + "user-triggered-fetchers-google.json"},
}

type Provider struct{}

func New() *Provider { return &Provider{} }

func (p *Provider) Meta() provider.Meta {
	return provider.Meta{
		ID:        "google",
		Name:      "Google",
		Homepage:  "https://developers.google.com/search/docs/crawling-indexing/verifying-googlebot",
		SourceURL: URLGoog,
		Category:  "cloud",
		Sets:      listSets(),
	}
}

func listSets() []provider.Set {
	out := make([]provider.Set, len(sets))
	for i, s := range sets {
		s.Where = map[string]string{"list": s.ID}
		out[i] = s
	}
	return out
}

func (p *Provider) Fetch(ctx context.Context, c *http.Client) ([]provider.Prefix, error) {
	var out []provider.Prefix
	for _, s := range sets {
		prefixes, err := provider.GetPrefixList(ctx, c, s.Source, map[string]string{"list": s.ID})
		if err != nil {
			return nil, err
		}
		out = append(out, prefixes...)
	}
	return out, nil
}
