package emit

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ipsets-io/ipsets/provider"
)

const SchemaVersion = 1

type Source struct {
	URL      string `json:"url"`
	Homepage string `json:"homepage,omitempty"`
}

type File struct {
	SchemaVersion int               `json:"schema_version"`
	Provider      string            `json:"provider"`
	ProviderName  string            `json:"provider_name"`
	Set           string            `json:"set"`
	SetName       string            `json:"set_name,omitempty"`
	Category      string            `json:"category,omitempty"`
	Family        string            `json:"family"`
	Source        Source            `json:"source"`
	Count         int               `json:"count"`
	Prefixes      []provider.Prefix `json:"prefixes"`
}

func (f File) Write(dir string) (string, error) {
	f.SchemaVersion = SchemaVersion
	f.Count = len(f.Prefixes)
	if f.Prefixes == nil {
		f.Prefixes = []provider.Prefix{}
	}

	body, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return "", err
	}
	body = append(body, '\n')
	if err := writeFile(filepath.Join(dir, f.Family+".json"), body); err != nil {
		return "", err
	}

	var txt bytes.Buffer
	for _, p := range f.Prefixes {
		txt.WriteString(p.Prefix.String())
		txt.WriteByte('\n')
	}
	if err := writeFile(filepath.Join(dir, f.Family+".txt"), txt.Bytes()); err != nil {
		return "", err
	}

	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:]), nil
}

func writeFile(path string, body []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, body, 0o644)
}
