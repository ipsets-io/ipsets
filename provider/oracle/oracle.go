package oracle

import (
	"context"
	"net/http"

	"github.com/ipsets-io/ipsets/provider"
)

const URL = "https://docs.oracle.com/en-us/iaas/tools/public_ip_ranges.json"

func tagged(id, name, tag string) provider.Set {
	return provider.Set{ID: id, Name: name, Where: map[string]string{"tag": tag}}
}

type Provider struct{}

func New() *Provider { return &Provider{} }

func (p *Provider) Meta() provider.Meta {
	return provider.Meta{
		ID:        "oracle",
		Name:      "Oracle Cloud Infrastructure",
		Homepage:  "https://docs.oracle.com/en-us/iaas/Content/General/Concepts/addressranges.htm",
		SourceURL: URL,
		Sets: []provider.Set{
			tagged("oci", "OCI customer infrastructure", "OCI"),
			tagged("osn", "Oracle Services Network", "OSN"),
			tagged("object-storage", "Object Storage", "OBJECT_STORAGE"),
		},
	}
}

type doc struct {
	Regions []struct {
		Region string `json:"region"`
		CIDRs  []struct {
			CIDR string   `json:"cidr"`
			Tags []string `json:"tags"`
		} `json:"cidrs"`
	} `json:"regions"`
}

func (p *Provider) Fetch(ctx context.Context, c *http.Client) ([]provider.Prefix, error) {
	var d doc
	if err := provider.GetJSON(ctx, c, URL, &d); err != nil {
		return nil, err
	}

	var out []provider.Prefix
	for _, r := range d.Regions {
		for _, e := range r.CIDRs {
			for _, tag := range e.Tags {
				parsed, err := provider.Parse([]string{e.CIDR}, map[string]string{
					"region": r.Region,
					"tag":    tag,
				})
				if err != nil {
					return nil, err
				}
				out = append(out, parsed...)
			}
		}
	}
	return out, nil
}
