package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/netip"
	"strings"
)

const UserAgent = "ipsets.io (+https://github.com/ipsets-io/ipsets)"

const maxBody = 64 << 20

func Get(ctx context.Context, c *http.Client, url string) ([]byte, error) {
	if c == nil {
		c = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBody))
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", url, err)
	}
	return body, nil
}

func GetJSON(ctx context.Context, c *http.Client, url string, v any) error {
	body, err := Get(ctx, c, url)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("decode %s: %w", url, err)
	}
	return nil
}

func GetLines(ctx context.Context, c *http.Client, url string) ([]string, error) {
	body, err := Get(ctx, c, url)
	if err != nil {
		return nil, err
	}
	var out []string
	for line := range strings.Lines(string(body)) {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			out = append(out, line)
		}
	}
	return out, nil
}

type prefixListDoc struct {
	Prefixes []struct {
		IPv4    string `json:"ipv4Prefix"`
		IPv6    string `json:"ipv6Prefix"`
		Service string `json:"service"`
		Scope   string `json:"scope"`
	} `json:"prefixes"`
}

func GetPrefixList(ctx context.Context, c *http.Client, url string, tags map[string]string) ([]Prefix, error) {
	var d prefixListDoc
	if err := GetJSON(ctx, c, url, &d); err != nil {
		return nil, err
	}

	out := make([]Prefix, 0, len(d.Prefixes))
	for _, e := range d.Prefixes {
		cidr := e.IPv4
		if cidr == "" {
			cidr = e.IPv6
		}
		if cidr == "" {
			continue
		}
		p, err := netip.ParsePrefix(strings.TrimSpace(cidr))
		if err != nil {
			return nil, fmt.Errorf("%s: parse %q: %w", url, cidr, err)
		}

		t := maps.Clone(tags)
		if t == nil {
			t = map[string]string{}
		}
		if e.Service != "" {
			t["service"] = e.Service
		}
		if e.Scope != "" {
			t["scope"] = e.Scope
		}
		out = append(out, Prefix{Prefix: p, Tags: t})
	}

	return out, nil
}

type prefixListProvider struct{ meta Meta }

func PrefixList(meta Meta) Provider { return prefixListProvider{meta: meta} }

func (p prefixListProvider) Meta() Meta { return p.meta }

func (p prefixListProvider) Fetch(ctx context.Context, c *http.Client) ([]Prefix, error) {
	return GetPrefixList(ctx, c, p.meta.SourceURL, nil)
}

func Parse(cidrs []string, tags map[string]string) ([]Prefix, error) {
	out := make([]Prefix, 0, len(cidrs))
	for _, s := range cidrs {
		s = strings.TrimSpace(s)
		p, err := netip.ParsePrefix(s)
		if err != nil {
			addr, addrErr := netip.ParseAddr(s)
			if addrErr != nil {
				return nil, fmt.Errorf("parse %q: %w", s, err)
			}
			p = netip.PrefixFrom(addr, addr.BitLen())
		}
		out = append(out, Prefix{Prefix: p, Tags: tags})
	}
	return out, nil
}
