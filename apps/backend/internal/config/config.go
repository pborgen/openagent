package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Framework string `json:"framework"`
	Workflow  string `json:"workflow"`
}

func Default() Config {
	return Config{Framework: "langgraph", Workflow: "default"}
}

// ConfigPath resolves repo-local config first, then ~/.openagent/config.json
func ConfigPath() string {
	if _, err := os.Stat("openagent.json"); err == nil {
		return "openagent.json"
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".openagent", "config.json")
}
