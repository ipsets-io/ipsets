package tor

import (
	"context"
	"net/http"

	"github.com/ipsets-io/ipsets/provider"
)

const URL = "https://check.torproject.org/torbulkexitlist"

type Provider struct{}

func New() *Provider { return &Provider{} }

func (p *Provider) Meta() provider.Meta {
	return provider.Meta{
		ID:        "tor",
		Name:      "Tor Project",
		Homepage:  "https://support.torproject.org/abuse/i-want-to-ban-tor/",
		SourceURL: URL,
		Sets: []provider.Set{
			{ID: "exits", Name: "Tor exit nodes", Category: "anonymizer"},
		},
	}
}

func (p *Provider) Fetch(ctx context.Context, c *http.Client) ([]provider.Prefix, error) {
	lines, err := provider.GetLines(ctx, c, URL)
	if err != nil {
		return nil, err
	}
	return provider.Parse(lines, nil)
}
