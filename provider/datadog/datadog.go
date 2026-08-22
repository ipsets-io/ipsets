package datadog

import (
	"context"
	"net/http"

	"github.com/ipsets-io/ipsets/provider"
)

const URL = "https://ip-ranges.datadoghq.com/"

type Provider struct{}

func New() *Provider { return &Provider{} }

func set(id, name, category string) provider.Set {
	return provider.Set{ID: id, Name: name, Category: category, Where: map[string]string{"service": id}}
}

func (p *Provider) Meta() provider.Meta {
	return provider.Meta{
		ID:        "datadog",
		Name:      "Datadog",
		Homepage:  "https://docs.datadoghq.com/api/latest/ip-ranges/",
		SourceURL: URL,
		Sets: []provider.Set{
			set("synthetics", "Synthetics probes", "monitoring"),
			set("webhooks", "Webhook senders", ""),
			set("agents", "Agents", ""),
			set("api", "API", ""),
			set("apm", "APM", ""),
			set("logs", "Logs", ""),
			set("process", "Process", ""),
			set("orchestrator", "Orchestrator", ""),
			set("remote-configuration", "Remote configuration", ""),
			set("synthetics-private-locations", "Synthetics private locations", ""),
			set("global", "Global", ""),
		},
	}
}

type block struct {
	V4 []string `json:"prefixes_ipv4"`
	V6 []string `json:"prefixes_ipv6"`
}

type doc struct {
	Synthetics        block `json:"synthetics"`
	Webhooks          block `json:"webhooks"`
	Agents            block `json:"agents"`
	API               block `json:"api"`
	APM               block `json:"apm"`
	Logs              block `json:"logs"`
	Process           block `json:"process"`
	Orchestrator      block `json:"orchestrator"`
	RemoteConfig      block `json:"remote-configuration"`
	SyntheticsPrivate block `json:"synthetics-private-locations"`
	Global            block `json:"global"`
}

func (p *Provider) Fetch(ctx context.Context, c *http.Client) ([]provider.Prefix, error) {
	var d doc
	if err := provider.GetJSON(ctx, c, URL, &d); err != nil {
		return nil, err
	}

	services := []struct {
		id string
		b  block
	}{
		{"synthetics", d.Synthetics},
		{"webhooks", d.Webhooks},
		{"agents", d.Agents},
		{"api", d.API},
		{"apm", d.APM},
		{"logs", d.Logs},
		{"process", d.Process},
		{"orchestrator", d.Orchestrator},
		{"remote-configuration", d.RemoteConfig},
		{"synthetics-private-locations", d.SyntheticsPrivate},
		{"global", d.Global},
	}

	var out []provider.Prefix
	for _, s := range services {
		parsed, err := provider.Parse(append(s.b.V4, s.b.V6...), map[string]string{"service": s.id})
		if err != nil {
			return nil, err
		}
		out = append(out, parsed...)
	}
	return out, nil
}
