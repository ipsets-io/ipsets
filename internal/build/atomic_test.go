package build

import (
	"context"
	"errors"
	"net/http"
	"net/netip"
	"os"
	"path/filepath"
	"testing"

	"github.com/ipsets-io/ipsets/provider"
)

type stub struct {
	id   string
	fail bool
}

func (s stub) Meta() provider.Meta {
	return provider.Meta{
		ID:        s.id,
		Name:      s.id,
		SourceURL: "https://example.test/list.json",
		Sets:      []provider.Set{{ID: "edge", Name: "Edge", Category: "cdn"}},
	}
}

func (s stub) Fetch(context.Context, *http.Client) ([]provider.Prefix, error) {
	if s.fail {
		return nil, errors.New("503 Service Unavailable")
	}
	return []provider.Prefix{{Prefix: netip.MustParsePrefix("1.2.3.0/24")}}, nil
}

func countFiles(t *testing.T, dir string) int {
	t.Helper()
	n := 0
	filepath.WalkDir(dir, func(_ string, d os.DirEntry, _ error) error {
		if d != nil && !d.IsDir() {
			n++
		}
		return nil
	})
	return n
}

func TestFailedBuildWritesNothing(t *testing.T) {
	dir := t.TempDir()

	// alpha sorts first and would have been written before bravo failed.
	providers := []provider.Provider{stub{id: "alpha"}, stub{id: "bravo", fail: true}}
	if _, err := Run(context.Background(), providers, Options{Dir: dir}); err == nil {
		t.Fatal("expected an error when a provider fails")
	}
	if n := countFiles(t, dir); n != 0 {
		t.Fatalf("a failed build must leave the published tree untouched, wrote %d files", n)
	}

	if _, err := Run(context.Background(), []provider.Provider{stub{id: "alpha"}}, Options{Dir: dir}); err != nil {
		t.Fatalf("healthy build: %v", err)
	}
	if countFiles(t, dir) == 0 {
		t.Fatal("a healthy build should write files")
	}
}

func TestBuildRemovesOrphanedSets(t *testing.T) {
	dir := t.TempDir()
	ghost := filepath.Join(dir, "v1", "ghost", "all")
	if err := os.MkdirAll(ghost, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ghost, "ipv4.txt"), []byte("1.2.3.0/24\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := Run(context.Background(), []provider.Provider{stub{id: "alpha"}}, Options{Dir: dir}); err != nil {
		t.Fatalf("build: %v", err)
	}

	if _, err := os.Stat(ghost); !os.IsNotExist(err) {
		t.Error("a set removed from the catalog must stop being served")
	}
	if _, err := os.Stat(filepath.Join(dir, "v1", "alpha", "edge", "ipv4.txt")); err != nil {
		t.Errorf("current sets must still be written: %v", err)
	}
}
