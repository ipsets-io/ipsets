package cloudflare

import (
	"context"
	"net/http"

	"github.com/ipsets-io/ipsets/provider"
)

const (
	URLv4 = "https://www.cloudflare.com/ips-v4/"
	URLv6 = "https://www.cloudflare.com/ips-v6/"
)

type Provider struct{}

func New() *Provider { return &Provider{} }

func (p *Provider) Meta() provider.Meta {
	return provider.Meta{
		ID:        "cloudflare",
		Name:      "Cloudflare",
		Homepage:  "https://www.cloudflare.com/ips/",
		SourceURL: URLv4,
		Sets: []provider.Set{
			{ID: "cdn", Name: "Cloudflare CDN", Category: "cdn"},
		},
	}
}

func (p *Provider) Fetch(ctx context.Context, c *http.Client) ([]provider.Prefix, error) {
	var out []provider.Prefix
	for _, url := range []string{URLv4, URLv6} {
		lines, err := provider.GetLines(ctx, c, url)
		if err != nil {
			return nil, err
		}
		prefixes, err := provider.Parse(lines, nil)
		if err != nil {
			return nil, err
		}
		out = append(out, prefixes...)
	}
	return out, nil
}
