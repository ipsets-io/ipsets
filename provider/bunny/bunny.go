package bunny

import (
	"context"
	"net/http"

	"github.com/ipsets-io/ipsets/provider"
)

const (
	URLv4 = "https://bunnycdn.com/api/system/edgeserverlist"
	URLv6 = "https://bunnycdn.com/api/system/edgeserverlist/ipv6"
)

type Provider struct{}

func New() *Provider { return &Provider{} }

func (p *Provider) Meta() provider.Meta {
	return provider.Meta{
		ID:        "bunny",
		Name:      "bunny.net",
		Homepage:  "https://bunny.net/",
		SourceURL: URLv4,
		Sets: []provider.Set{
			{ID: "cdn", Name: "bunny.net edge servers", Category: "cdn"},
		},
	}
}

func (p *Provider) Fetch(ctx context.Context, c *http.Client) ([]provider.Prefix, error) {
	var out []provider.Prefix
	for _, url := range []string{URLv4, URLv6} {
		var addrs []string
		if err := provider.GetJSON(ctx, c, url, &addrs); err != nil {
			return nil, err
		}
		parsed, err := provider.Parse(addrs, nil)
		if err != nil {
			return nil, err
		}
		out = append(out, parsed...)
	}
	return out, nil
}
