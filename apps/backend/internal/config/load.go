package config

import (
	"encoding/json"
	"os"
)

func Load() (Config, error) {
	cfg := Default()
	path := ConfigPath()
	b, err := os.ReadFile(path)
	if err != nil {
		return cfg, nil // default if missing
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
