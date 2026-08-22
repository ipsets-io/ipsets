package pingdom

import (
	"context"
	"net/http"

	"github.com/ipsets-io/ipsets/provider"
)

const (
	URLv4 = "https://my.pingdom.com/probes/ipv4"
	URLv6 = "https://my.pingdom.com/probes/ipv6"
)

type Provider struct{}

func New() *Provider { return &Provider{} }

func (p *Provider) Meta() provider.Meta {
	return provider.Meta{
		ID:        "pingdom",
		Name:      "Pingdom",
		Homepage:  "https://www.pingdom.com/rss/probe_servers.xml",
		SourceURL: URLv4,
		Sets: []provider.Set{
			{ID: "probes", Name: "Pingdom probe servers", Category: "monitoring"},
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
		parsed, err := provider.Parse(lines, nil)
		if err != nil {
			return nil, err
		}
		out = append(out, parsed...)
	}
	return out, nil
}
