package provider

import (
	"context"
	"net/http"
	"net/netip"
)

type Prefix struct {
	Prefix netip.Prefix      `json:"prefix"`
	Tags   map[string]string `json:"tags,omitempty"`
}

type Set struct {
	ID       string
	Name     string
	Category string
	Source   string
	Where    map[string]string
}

type Meta struct {
	ID        string
	Name      string
	Homepage  string
	SourceURL string
	Sets      []Set
}

type Provider interface {
	Meta() Meta
	Fetch(ctx context.Context, c *http.Client) ([]Prefix, error)
}

func (s Set) Matches(p Prefix) bool {
	for k, want := range s.Where {
		if p.Tags[k] != want {
			return false
		}
	}
	return true
}
