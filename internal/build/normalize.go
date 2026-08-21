package build

import (
	"cmp"
	"net/netip"
	"slices"

	"github.com/ipsets-io/ipsets/provider"
)

func Split(in []provider.Prefix) (v4, v6 []provider.Prefix) {
	for _, p := range in {
		if !p.Prefix.IsValid() {
			continue
		}
		if p.Prefix.Addr().Is4() {
			v4 = append(v4, p)
		} else {
			v6 = append(v6, p)
		}
	}
	return v4, v6
}

func Normalize(in []provider.Prefix) []provider.Prefix {
	byPrefix := make(map[netip.Prefix]map[string]string, len(in))
	order := make([]netip.Prefix, 0, len(in))

	for _, p := range in {
		masked := p.Prefix.Masked()
		tags := cleanTags(p.Tags)
		existing, seen := byPrefix[masked]
		if !seen {
			byPrefix[masked] = tags
			order = append(order, masked)
			continue
		}
		byPrefix[masked] = intersectTags(existing, tags)
	}

	out := make([]provider.Prefix, 0, len(order))
	for _, p := range order {
		tags := byPrefix[p]
		if len(tags) == 0 {
			tags = nil
		}
		out = append(out, provider.Prefix{Prefix: p, Tags: tags})
	}

	slices.SortFunc(out, func(a, b provider.Prefix) int {
		if c := a.Prefix.Addr().Compare(b.Prefix.Addr()); c != 0 {
			return c
		}
		return cmp.Compare(a.Prefix.Bits(), b.Prefix.Bits())
	})
	return out
}

func cleanTags(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		if v != "" {
			out[k] = v
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func intersectTags(a, b map[string]string) map[string]string {
	out := make(map[string]string, min(len(a), len(b)))
	for k, v := range a {
		if b[k] == v {
			out[k] = v
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func Collapse(sorted []provider.Prefix) []provider.Prefix {
	out := sorted[:0:0]
	for _, p := range sorted {
		if n := len(out); n > 0 {
			container := out[n-1].Prefix
			if container.Overlaps(p.Prefix) && p.Prefix.Bits() >= container.Bits() {
				out[n-1].Tags = intersectTags(out[n-1].Tags, p.Tags)
				continue
			}
		}
		out = append(out, p)
	}
	return out
}
