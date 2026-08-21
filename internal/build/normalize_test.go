package build

import (
	"net/netip"
	"testing"

	"github.com/ipsets-io/ipsets/provider"
)

func prefixes(t *testing.T, in ...string) []provider.Prefix {
	t.Helper()
	out := make([]provider.Prefix, 0, len(in))
	for _, s := range in {
		out = append(out, provider.Prefix{Prefix: netip.MustParsePrefix(s)})
	}
	return out
}

func strs(in []provider.Prefix) []string {
	out := make([]string, 0, len(in))
	for _, p := range in {
		out = append(out, p.Prefix.String())
	}
	return out
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestNormalizeMasksSortsAndDedupes(t *testing.T) {
	in := prefixes(t, "10.0.0.0/8", "1.2.3.4/24", "10.0.0.0/8", "1.2.3.0/24")
	got := strs(Normalize(in))
	want := []string{"1.2.3.0/24", "10.0.0.0/8"}
	if !equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestNormalizeKeepsOnlyUnanimousTags(t *testing.T) {
	p := netip.MustParsePrefix("23.228.249.0/24")
	in := []provider.Prefix{
		{Prefix: p, Tags: map[string]string{"region": "GLOBAL", "service": "AMAZON"}},
		{Prefix: p, Tags: map[string]string{"region": "GLOBAL", "service": "CLOUDFRONT"}},
	}

	got := Normalize(in)
	if len(got) != 1 {
		t.Fatalf("expected 1 prefix, got %d", len(got))
	}
	if got[0].Tags["region"] != "GLOBAL" {
		t.Errorf("region should survive, got %q", got[0].Tags["region"])
	}
	if _, ok := got[0].Tags["service"]; ok {
		t.Errorf("conflicting service tag should be dropped, got %q", got[0].Tags["service"])
	}
}

func TestNormalizeDropsEmptyTagValues(t *testing.T) {
	in := []provider.Prefix{{
		Prefix: netip.MustParsePrefix("1.2.3.0/24"),
		Tags:   map[string]string{"region": "", "service": "EC2"},
	}}

	got := Normalize(in)
	if _, ok := got[0].Tags["region"]; ok {
		t.Errorf("empty tag value should be dropped: %v", got[0].Tags)
	}
	if got[0].Tags["service"] != "EC2" {
		t.Errorf("service tag should survive: %v", got[0].Tags)
	}
}

func TestSplitByFamily(t *testing.T) {
	v4, v6 := Split(prefixes(t, "1.2.3.0/24", "2001:db8::/32", "10.0.0.0/8"))
	if len(v4) != 2 || len(v6) != 1 {
		t.Fatalf("got v4=%d v6=%d, want 2 and 1", len(v4), len(v6))
	}
}

func TestCollapseRemovesContainedPrefixes(t *testing.T) {
	in := Normalize(prefixes(t,
		"130.176.0.0/18", "130.176.128.0/21", "130.176.0.0/24", "131.0.0.0/16"))
	got := strs(Collapse(in))
	want := []string{"130.176.0.0/18", "130.176.128.0/21", "131.0.0.0/16"}
	if !equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestCollapseMergesTagsOfAbsorbedPrefixes(t *testing.T) {
	in := Normalize([]provider.Prefix{
		{Prefix: netip.MustParsePrefix("10.0.0.0/8"), Tags: map[string]string{"service": "EC2", "region": "us-east-1"}},
		{Prefix: netip.MustParsePrefix("10.5.0.0/16"), Tags: map[string]string{"service": "S3", "region": "us-east-1"}},
	})

	got := Collapse(in)
	if len(got) != 1 {
		t.Fatalf("expected 1 prefix, got %d", len(got))
	}
	if _, ok := got[0].Tags["service"]; ok {
		t.Errorf("absorbing a prefix with a different service must drop the tag, got %q", got[0].Tags["service"])
	}
	if got[0].Tags["region"] != "us-east-1" {
		t.Errorf("agreed tag should survive, got %v", got[0].Tags)
	}
}

func TestCollapseLeavesNoContainedPrefixes(t *testing.T) {
	in := prefixes(t,
		"10.5.1.0/24", "10.0.0.0/8", "10.5.0.0/16", "11.0.0.0/8",
		"192.168.1.0/24", "192.168.0.0/16", "172.16.0.0/12", "172.16.5.0/24",
		"130.176.128.0/21", "130.176.0.0/18", "8.8.8.8/32", "8.8.8.0/24",
		"2001:db8::/32", "2001:db8:1::/48", "2001:db9::/32",
	)

	v4, v6 := Split(in)
	for _, family := range [][]provider.Prefix{v4, v6} {
		got := Collapse(Normalize(family))
		for i, a := range got {
			for j, b := range got {
				if i != j && b.Prefix.Overlaps(a.Prefix) && b.Prefix.Bits() > a.Prefix.Bits() {
					t.Errorf("%s still contains %s after collapse", a.Prefix, b.Prefix)
				}
			}
		}
	}
}

func TestCollapseKeepsDisjointPrefixes(t *testing.T) {
	in := Normalize(prefixes(t, "10.0.0.0/9", "10.128.0.0/9"))
	got := strs(Collapse(in))
	want := []string{"10.0.0.0/9", "10.128.0.0/9"}
	if !equal(got, want) {
		t.Fatalf("adjacent siblings must not be merged: got %v, want %v", got, want)
	}
}
