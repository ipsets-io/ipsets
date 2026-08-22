package azure

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ipsets-io/ipsets/provider"
)

const (
	urlPrefix = "https://download.microsoft.com/download/7/1/d/71d86715-5596-4529-9b13-da13a5de5b63/ServiceTags_Public_"
	downloads = "https://www.microsoft.com/en-us/download/details.aspx?id=56519"
)

const lookbackDays = 21

func tagged(id, name, category, tag string) provider.Set {
	return provider.Set{ID: id, Name: name, Category: category, Where: map[string]string{"service": tag}}
}

type Provider struct{}

func New() *Provider { return &Provider{} }

func (p *Provider) Meta() provider.Meta {
	return provider.Meta{
		ID:        "azure",
		Name:      "Microsoft Azure",
		Homepage:  "https://learn.microsoft.com/en-us/azure/virtual-network/service-tags-overview",
		SourceURL: downloads,
		Sets: []provider.Set{
			tagged("front-door-backend", "Front Door (origin-facing)", "cdn", "AzureFrontDoor.Backend"),
			tagged("front-door-frontend", "Front Door (edge)", "cdn", "AzureFrontDoor.Frontend"),
			tagged("devops", "Azure DevOps hosted agents", "ci", "AzureDevOps"),
			tagged("application-insights-availability", "Application Insights availability tests", "monitoring", "ApplicationInsightsAvailability"),
			tagged("traffic-manager", "Traffic Manager probes", "monitoring", "AzureTrafficManager"),
			tagged("monitor", "Azure Monitor", "", "AzureMonitor"),
			tagged("app-service", "App Service", "", "AppService"),
			tagged("storage", "Storage", "", "Storage"),
			tagged("sql", "SQL", "", "Sql"),
			tagged("active-directory", "Entra ID", "", "AzureActiveDirectory"),
			tagged("container-registry", "Container Registry", "", "AzureContainerRegistry"),
		},
	}
}

type doc struct {
	Values []struct {
		Name       string `json:"name"`
		Properties struct {
			Region          string   `json:"region"`
			AddressPrefixes []string `json:"addressPrefixes"`
		} `json:"properties"`
	} `json:"values"`
}

// Microsoft publishes weekly under a dated filename and keeps roughly two weeks
// online, so the current file has to be discovered by trying recent dates.
func fetch(ctx context.Context, c *http.Client) (doc, error) {
	var lastErr error
	for i := range lookbackDays {
		day := time.Now().UTC().AddDate(0, 0, -i).Format("20060102")
		var d doc
		if err := provider.GetJSON(ctx, c, urlPrefix+day+".json", &d); err != nil {
			lastErr = err
			continue
		}
		return d, nil
	}
	return doc{}, fmt.Errorf("no service tag file found in the last %d days: %w", lookbackDays, lastErr)
}

func (p *Provider) Fetch(ctx context.Context, c *http.Client) ([]provider.Prefix, error) {
	d, err := fetch(ctx, c)
	if err != nil {
		return nil, err
	}

	// Regional variants (Storage.WestEurope) list the same prefix strings as their
	// base tag, so they are the only way to recover a region for a base-tag prefix.
	region := map[string]string{}
	for _, v := range d.Values {
		r := v.Properties.Region
		if r == "" {
			continue
		}
		for _, p := range v.Properties.AddressPrefixes {
			if prev, ok := region[p]; ok && prev != r {
				region[p] = ""
				continue
			}
			region[p] = r
		}
	}

	var out []provider.Prefix
	for _, v := range d.Values {
		if v.Properties.Region != "" {
			continue
		}
		for _, p := range v.Properties.AddressPrefixes {
			parsed, err := provider.Parse([]string{p}, map[string]string{
				"service": v.Name,
				"region":  region[p],
			})
			if err != nil {
				return nil, err
			}
			out = append(out, parsed...)
		}
	}
	return out, nil
}
