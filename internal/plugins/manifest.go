package plugins

import (
	"encoding/json"
	"os"
	"path"
)

const CURRENT_MANIFEST_VERSION uint = 2

type ManifestPlatform struct {
	Os   []string
	Arch []string
}

type Manifest struct {
	ManifestVersion uint             `json:"manifest_version"`
	Name            string           `json:"name"`
	PluginVersion   []uint           `json:"plugin_version"`
	Description     string           `json:"description"`
	Author          []string         `json:"author"`
	Entry           string           `json:"entry"`
	Platform        ManifestPlatform `json:"platform"`
}

func readManifestTo(pluginPath string, m *Manifest) error {
	b, err := os.ReadFile(path.Join(pluginPath, "manifest.json"))
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, m)
	if err != nil {
		return err
	}
	return nil
}
