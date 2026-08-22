package emit

import (
	"bytes"
	_ "embed"
	"html/template"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

//go:embed index.html.tmpl
var siteTemplate string

type SiteSet struct {
	ID       string
	Name     string
	Category string
	Path     string
	Counts   string
}

type SiteGroup struct {
	ID       string
	Name     string
	Homepage string
	Sets     []SiteSet
}

type SiteData struct {
	Providers  []SiteGroup
	Categories []SiteSet
	Built      string
	Counts     struct {
		Providers int
		Sets      int
	}
}

func (i Index) Site() SiteData {
	var d SiteData

	for _, p := range i.Providers {
		g := SiteGroup{ID: p.ID, Name: p.Name, Homepage: p.Homepage}
		for _, s := range p.Sets {
			g.Sets = append(g.Sets, SiteSet{
				ID:       s.ID,
				Name:     s.Name,
				Category: s.Category,
				Path:     "/" + path.Join("v1", p.ID, s.ID),
				Counts:   countLabel(s),
			})
		}
		d.Counts.Sets += len(g.Sets)
		d.Providers = append(d.Providers, g)
	}
	d.Counts.Providers = len(d.Providers)

	for _, c := range i.Categories {
		set := SiteSet{ID: c.ID, Name: c.ID, Category: c.ID,
			Path: "/" + path.Join("v1", "categories", c.ID)}
		if len(c.Sets) > 0 {
			set.Counts = countLabel(c.Sets[0])
		}
		d.Categories = append(d.Categories, set)
	}

	if !i.GeneratedAt.IsZero() {
		d.Built = i.GeneratedAt.Format("2 January 2006")
	}
	return d
}

func countLabel(s IndexSet) string {
	var parts []string
	if s.IPv4.Count > 0 {
		parts = append(parts, thousands(s.IPv4.Count)+" v4")
	}
	if s.IPv6.Count > 0 {
		parts = append(parts, thousands(s.IPv6.Count)+" v6")
	}
	return strings.Join(parts, " · ")
}

func thousands(n int) string {
	s := strconv.Itoa(n)
	if len(s) <= 3 {
		return s
	}
	var b strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(c)
	}
	return b.String()
}

const Domain = "ipsets.io"

func (i Index) WriteSite(dir string) error {
	tmpl, err := template.New("index").Parse(siteTemplate)
	if err != nil {
		return err
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, i.Site()); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(dir, "index.html"), b.Bytes()); err != nil {
		return err
	}

	return writeFile(filepath.Join(dir, "CNAME"), []byte(Domain+"\n"))
}
