package sans

import (
	"context"
	"fmt"
	"net/http"
	"net/netip"
	"strings"

	"github.com/ipsets-io/ipsets/provider"
)

// URL is DShield's (SANS Internet Storm Center) recommended block list: the
// top 20 attacking /24 subnets by distinct targets scanned over the last
// three days. Unlike every other provider here, this is not a vendor
// declaring its own infrastructure — it's third-party abuse telemetry that
// churns every few days, so treat it as reputation data, not attribution.
const URL = "https://isc.sans.edu/block.txt"

type Provider struct{}

func New() *Provider { return &Provider{} }

func (p *Provider) Meta() provider.Meta {
	return provider.Meta{
		ID:        "sans",
		Name:      "SANS Internet Storm Center",
		Homepage:  "https://isc.sans.edu/",
		SourceURL: URL,
		Sets: []provider.Set{
			{ID: "top-attackers", Name: "Top attacking /24 subnets (DShield)", Category: "abuse"},
		},
	}
}

func (p *Provider) Fetch(ctx context.Context, c *http.Client) ([]provider.Prefix, error) {
	lines, err := provider.GetLines(ctx, c, URL)
	if err != nil {
		return nil, err
	}

	out := make([]provider.Prefix, 0, len(lines))
	for _, line := range lines {
		fields := strings.Split(line, "\t")
		if len(fields) < 7 {
			return nil, fmt.Errorf("%s: unexpected line %q", URL, line)
		}

		start, maskLen, attacks, name, country, contact := fields[0], fields[2], fields[3], fields[4], fields[5], fields[6]

		p, err := netip.ParsePrefix(fmt.Sprintf("%s/%s", start, maskLen))
		if err != nil {
			return nil, fmt.Errorf("%s: parse %q/%q: %w", URL, start, maskLen, err)
		}

		tags := map[string]string{"attacks": attacks}
		if name != "-" {
			tags["network"] = name
		}
		if country != "-" {
			tags["country"] = country
		}
		if contact != "-" && contact != "None" && contact != ">>UNKNOWN<<" {
			tags["contact"] = contact
		}

		out = append(out, provider.Prefix{Prefix: p, Tags: tags})
	}

	return out, nil
}
