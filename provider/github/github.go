package github

import (
	"context"
	"net/http"

	"github.com/ipsets-io/ipsets/provider"
)

const URL = "https://api.github.com/meta"

type Provider struct{}

func New() *Provider { return &Provider{} }

func usage(id, name, category string) provider.Set {
	return provider.Set{
		ID:       id,
		Name:     name,
		Category: category,
		Where:    map[string]string{"usage": id},
	}
}

func (p *Provider) Meta() provider.Meta {
	return provider.Meta{
		ID:        "github",
		Name:      "GitHub",
		Homepage:  "https://docs.github.com/en/rest/meta/meta",
		SourceURL: URL,
		Sets: []provider.Set{
			usage("actions", "Actions", "ci"),
			usage("actions-macos", "Actions (macOS)", "ci"),
			usage("codespaces", "Codespaces", "ci"),
			usage("hooks", "Webhooks", ""),
			usage("web", "Web", ""),
			usage("api", "API", ""),
			usage("git", "Git", ""),
			usage("packages", "Packages", ""),
			usage("pages", "Pages", ""),
			usage("importer", "Importer", ""),
			usage("enterprise-importer", "Enterprise Importer", ""),
			usage("copilot", "Copilot", ""),
		},
	}
}

type doc struct {
	Hooks              []string `json:"hooks"`
	Web                []string `json:"web"`
	API                []string `json:"api"`
	Git                []string `json:"git"`
	Packages           []string `json:"packages"`
	Pages              []string `json:"pages"`
	Importer           []string `json:"importer"`
	EnterpriseImporter []string `json:"github_enterprise_importer"`
	Actions            []string `json:"actions"`
	ActionsMacos       []string `json:"actions_macos"`
	Codespaces         []string `json:"codespaces"`
	Copilot            []string `json:"copilot"`
}

func (p *Provider) Fetch(ctx context.Context, c *http.Client) ([]provider.Prefix, error) {
	var d doc
	if err := provider.GetJSON(ctx, c, URL, &d); err != nil {
		return nil, err
	}

	lists := []struct {
		usage string
		cidrs []string
	}{
		{"hooks", d.Hooks},
		{"web", d.Web},
		{"api", d.API},
		{"git", d.Git},
		{"packages", d.Packages},
		{"pages", d.Pages},
		{"importer", d.Importer},
		{"enterprise-importer", d.EnterpriseImporter},
		{"actions", d.Actions},
		{"actions-macos", d.ActionsMacos},
		{"codespaces", d.Codespaces},
		{"copilot", d.Copilot},
	}

	var out []provider.Prefix
	for _, l := range lists {
		parsed, err := provider.Parse(l.cidrs, map[string]string{"usage": l.usage})
		if err != nil {
			return nil, err
		}
		out = append(out, parsed...)
	}
	return out, nil
}
