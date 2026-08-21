package build

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"path"
	"path/filepath"
	"slices"
	"time"

	"github.com/ipsets-io/ipsets/internal/emit"
	"github.com/ipsets-io/ipsets/provider"
)

const (
	APIVersion       = "v1"
	CategoryProvider = "categories"
)

type Options struct {
	Dir    string
	Client *http.Client
	Now    time.Time
}

func (o *Options) setDefaults() {
	if o.Client == nil {
		o.Client = &http.Client{Timeout: 90 * time.Second}
	}
	if o.Now.IsZero() {
		o.Now = time.Now()
	}
	o.Now = o.Now.UTC().Truncate(time.Second)
}

type fetched struct {
	meta     provider.Meta
	prefixes []provider.Prefix
}

func Run(ctx context.Context, providers []provider.Provider, opts Options) (emit.Index, error) {
	opts.setDefaults()

	all := make([]fetched, 0, len(providers))
	for _, p := range providers {
		meta := p.Meta()
		prefixes, err := p.Fetch(ctx, opts.Client)
		if err != nil {
			return emit.Index{}, fmt.Errorf("%s: %w", meta.ID, err)
		}
		all = append(all, fetched{meta: meta, prefixes: prefixes})
	}

	root := filepath.Join(opts.Dir, APIVersion)
	index := emit.Index{GeneratedAt: opts.Now}
	categories := map[string][]provider.Prefix{}

	for _, f := range all {
		entry := emit.IndexProvider{
			ID:       f.meta.ID,
			Name:     f.meta.Name,
			Homepage: f.meta.Homepage,
			Source:   f.meta.SourceURL,
		}

		for _, set := range f.meta.AllSets() {
			selected := selectPrefixes(f.prefixes, set)
			if len(selected) == 0 {
				return emit.Index{}, fmt.Errorf("%s/%s: matched no prefixes", f.meta.ID, set.ID)
			}

			out, err := writeSet(root, f.meta, set, selected)
			if err != nil {
				return emit.Index{}, fmt.Errorf("%s/%s: %w", f.meta.ID, set.ID, err)
			}

			entry.Sets = append(entry.Sets, out)
			if set.Category != "" {
				categories[set.Category] = append(categories[set.Category], tagProvider(selected, f.meta.ID)...)
			}
		}
		index.Providers = append(index.Providers, entry)
	}

	cats, err := writeCategories(root, categories)
	if err != nil {
		return emit.Index{}, err
	}
	index.Categories = cats

	if err := index.Write(root); err != nil {
		return emit.Index{}, err
	}
	if err := index.WriteSite(opts.Dir); err != nil {
		return emit.Index{}, err
	}
	return index, nil
}

func selectPrefixes(all []provider.Prefix, set provider.Set) []provider.Prefix {
	if len(set.Where) == 0 {
		return all
	}
	out := make([]provider.Prefix, 0, len(all))
	for _, p := range all {
		if set.Matches(p) {
			out = append(out, p)
		}
	}
	return out
}

func tagProvider(in []provider.Prefix, id string) []provider.Prefix {
	out := make([]provider.Prefix, len(in))
	for i, p := range in {
		tags := maps.Clone(p.Tags)
		if tags == nil {
			tags = map[string]string{}
		}
		tags["provider"] = id
		out[i] = provider.Prefix{Prefix: p.Prefix, Tags: tags}
	}
	return out
}

func writeSet(root string, meta provider.Meta, set provider.Set, selected []provider.Prefix) (emit.IndexSet, error) {
	dir := filepath.Join(root, meta.ID, set.ID)

	source := set.Source
	if source == "" {
		source = meta.SourceURL
	}
	file := emit.File{
		Provider:     meta.ID,
		ProviderName: meta.Name,
		Set:          set.ID,
		SetName:      set.Name,
		Category:     set.Category,
		Source:       emit.Source{URL: source, Homepage: meta.Homepage},
	}

	v4, v6 := Split(selected)
	out := emit.IndexSet{ID: set.ID, Name: set.Name, Category: set.Category}

	var err error
	if out.IPv4, err = writeFamily(dir, "ipv4", prepare(v4), file, meta.ID, set.ID); err != nil {
		return out, err
	}
	if out.IPv6, err = writeFamily(dir, "ipv6", prepare(v6), file, meta.ID, set.ID); err != nil {
		return out, err
	}
	return out, nil
}

func writeFamily(dir, family string, prefixes []provider.Prefix, file emit.File, providerID, setID string) (emit.IndexFile, error) {
	for _, p := range prefixes {
		if p.Prefix.Bits() == 0 {
			return emit.IndexFile{}, fmt.Errorf("refusing to publish default route %s", p.Prefix)
		}
	}

	file.Family = family
	file.Prefixes = prefixes

	sum, err := file.Write(dir)
	if err != nil {
		return emit.IndexFile{}, err
	}
	return emit.IndexFile{
		Path:   "/" + path.Join(APIVersion, providerID, setID, family+".json"),
		Count:  len(prefixes),
		SHA256: sum,
	}, nil
}

func prepare(in []provider.Prefix) []provider.Prefix {
	return Collapse(Normalize(in))
}

func writeCategories(root string, categories map[string][]provider.Prefix) ([]emit.IndexCategory, error) {
	out := make([]emit.IndexCategory, 0, len(categories))

	for _, id := range slices.Sorted(maps.Keys(categories)) {
		meta := provider.Meta{ID: CategoryProvider, Name: "Categories"}
		set := provider.Set{ID: id, Name: id, Category: id}

		entry, err := writeSet(root, meta, set, categories[id])
		if err != nil {
			return nil, fmt.Errorf("categories/%s: %w", id, err)
		}
		out = append(out, emit.IndexCategory{ID: id, Sets: []emit.IndexSet{entry}})
	}
	return out, nil
}
