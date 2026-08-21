package fastly

import (
	"context"
	"net/http"

	"github.com/ipsets-io/ipsets/provider"
)

const URL = "https://api.fastly.com/public-ip-list"

type Provider struct{}

func New() *Provider { return &Provider{} }

func (p *Provider) Meta() provider.Meta {
	return provider.Meta{
		ID:        "fastly",
		Name:      "Fastly",
		Homepage:  "https://developer.fastly.com/reference/api/utils/public-ip-list/",
		SourceURL: URL,
		Category:  "cdn",
	}
}

type doc struct {
	Addresses     []string `json:"addresses"`
	IPv6Addresses []string `json:"ipv6_addresses"`
}

func (p *Provider) Fetch(ctx context.Context, c *http.Client) ([]provider.Prefix, error) {
	var d doc
	if err := provider.GetJSON(ctx, c, URL, &d); err != nil {
		return nil, err
	}
	return provider.Parse(append(d.Addresses, d.IPv6Addresses...), nil)
}
