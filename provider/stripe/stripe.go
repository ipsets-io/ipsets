package stripe

import (
	"context"
	"maps"
	"net/http"
	"slices"

	"github.com/ipsets-io/ipsets/provider"
)

const base = "https://stripe.com/files/ips/ips_"

var sets = []provider.Set{
	{ID: "api", Name: "Stripe API", Source: base + "api.json"},
	{ID: "armada-gator", Name: "Files, Armada and Gator", Source: base + "armada_gator.json"},
	{ID: "webhooks", Name: "Webhook senders", Source: base + "webhooks.json"},
}

type Provider struct{}

func New() *Provider { return &Provider{} }

func (p *Provider) Meta() provider.Meta {
	return provider.Meta{
		ID:        "stripe",
		Name:      "Stripe",
		Homepage:  "https://docs.stripe.com/ips",
		SourceURL: sets[0].Source,
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
		var d map[string][]string
		if err := provider.GetJSON(ctx, c, s.Source, &d); err != nil {
			return nil, err
		}
		for _, k := range slices.Sorted(maps.Keys(d)) {
			parsed, err := provider.Parse(d[k], map[string]string{"list": s.ID})
			if err != nil {
				return nil, err
			}
			out = append(out, parsed...)
		}
	}
	return out, nil
}
