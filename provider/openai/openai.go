package openai

import (
	"context"
	"net/http"

	"github.com/ipsets-io/ipsets/provider"
)

const URLGPTBot = "https://openai.com/gptbot.json"

var sets = []provider.Set{
	{ID: "gptbot", Name: "GPTBot (model training)", Category: "crawler", Source: URLGPTBot},
	{ID: "chatgpt-user", Name: "ChatGPT-User (user-triggered fetches)", Category: "crawler", Source: "https://openai.com/chatgpt-user.json"},
	{ID: "searchbot", Name: "OAI-SearchBot", Category: "crawler", Source: "https://openai.com/searchbot.json"},
}

type Provider struct{}

func New() *Provider { return &Provider{} }

func (p *Provider) Meta() provider.Meta {
	out := make([]provider.Set, len(sets))
	for i, s := range sets {
		s.Where = map[string]string{"list": s.ID}
		out[i] = s
	}
	return provider.Meta{
		ID:        "openai",
		Name:      "OpenAI",
		Homepage:  "https://platform.openai.com/docs/bots",
		SourceURL: URLGPTBot,
		Sets:      out,
	}
}

func (p *Provider) Fetch(ctx context.Context, c *http.Client) ([]provider.Prefix, error) {
	var out []provider.Prefix
	for _, s := range sets {
		prefixes, err := provider.GetPrefixList(ctx, c, s.Source, map[string]string{"list": s.ID})
		if err != nil {
			return nil, err
		}
		out = append(out, prefixes...)
	}
	return out, nil
}
