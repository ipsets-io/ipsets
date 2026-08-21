package emit

import (
	"encoding/json"
	"path/filepath"
	"time"
)

type Index struct {
	SchemaVersion int             `json:"schema_version"`
	GeneratedAt   time.Time       `json:"generated_at"`
	Providers     []IndexProvider `json:"providers"`
	Categories    []IndexCategory `json:"categories"`
}

type IndexProvider struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Homepage string     `json:"homepage,omitempty"`
	Source   string     `json:"source,omitempty"`
	Sets     []IndexSet `json:"sets"`
}

type IndexCategory struct {
	ID   string     `json:"id"`
	Sets []IndexSet `json:"sets"`
}

type IndexSet struct {
	ID       string    `json:"id"`
	Name     string    `json:"name,omitempty"`
	Category string    `json:"category,omitempty"`
	IPv4     IndexFile `json:"ipv4"`
	IPv6     IndexFile `json:"ipv6"`
}

type IndexFile struct {
	Path   string `json:"path"`
	Count  int    `json:"count"`
	SHA256 string `json:"sha256,omitempty"`
}

func (i Index) Write(dir string) error {
	i.SchemaVersion = SchemaVersion
	body, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err
	}
	return writeFile(filepath.Join(dir, "index.json"), append(body, '\n'))
}
