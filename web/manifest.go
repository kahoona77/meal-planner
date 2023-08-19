package web

import (
	"encoding/json"
	"io"
	"io/fs"
)

type Manifest map[string]ManifestEntry

type ManifestEntry struct {
	File string `json:"file"`
	Src  string `json:"src"`
}

func (m Manifest) File(key string) string {
	entry, ok := m[key]
	if !ok {
		return key
	}
	return entry.File
}

func EmptyManifest() Manifest {
	return make(Manifest, 0)
}

func ParseManifest(file string, fSys fs.FS) (Manifest, error) {
	f, err := fSys.Open(file)
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	err = json.Unmarshal(content, &manifest)
	return manifest, err
}
