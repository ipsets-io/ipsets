package sentry

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"slices"

	"github.com/ipsets-io/ipsets/provider"
)

const URL = "https://docs.sentry.io/api/ip-ranges/"

type Provider struct{}

func New() *Provider { return &Provider{} }

func set(id, name, category, service string) provider.Set {
	return provider.Set{ID: id, Name: name, Category: category, Where: map[string]string{"service": service}}
}

func (p *Provider) Meta() provider.Meta {
	return provider.Meta{
		ID:        "sentry",
		Name:      "Sentry",
		Homepage:  "https://docs.sentry.io/security-legal-pii/security/ip-ranges/",
		SourceURL: URL,
		Sets: []provider.Set{
			set("uptime-monitoring", "Uptime monitoring probes", "monitoring", "uptime_monitoring"),
			set("outbound-requests", "Outbound requests (webhooks, integrations)", "", "outbound_requests"),
			set("email-delivery", "Email delivery", "", "email_delivery"),
			set("event-ingestion", "Event ingestion endpoints", "", "event_ingestion"),
			set("dashboard", "Dashboard", "", "dashboard"),
		},
	}
}

type doc struct {
	Data map[string]any `json:"data"`
}

// Sentry nests irregularly: a service maps either to a list of prefixes or to
// further named groups, sometimes two levels deep. The path below the service
// name is kept as a scope tag rather than modelled.
func collect(node any, service, scope string) ([]provider.Prefix, error) {
	switch v := node.(type) {
	case []any:
		cidrs := make([]string, 0, len(v))
		for _, item := range v {
			s, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("%s: expected a prefix string, got %T", service, item)
			}
			cidrs = append(cidrs, s)
		}
		return provider.Parse(cidrs, map[string]string{"service": service, "scope": scope})

	case map[string]any:
		var out []provider.Prefix
		for _, k := range slices.Sorted(maps.Keys(v)) {
			child := k
			if scope != "" {
				child = scope + "." + k
			}
			got, err := collect(v[k], service, child)
			if err != nil {
				return nil, err
			}
			out = append(out, got...)
		}
		return out, nil
	}
	return nil, fmt.Errorf("%s: unexpected value of type %T", service, node)
}

func (p *Provider) Fetch(ctx context.Context, c *http.Client) ([]provider.Prefix, error) {
	var d doc
	if err := provider.GetJSON(ctx, c, URL, &d); err != nil {
		return nil, err
	}

	var out []provider.Prefix
	for _, service := range slices.Sorted(maps.Keys(d.Data)) {
		got, err := collect(d.Data[service], service, "")
		if err != nil {
			return nil, err
		}
		out = append(out, got...)
	}
	return out, nil
}
