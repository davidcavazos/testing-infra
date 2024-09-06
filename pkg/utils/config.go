package utils

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	Match   []string `json:"match"`
	Ignore  []string `json:"ignore"`
	Package []string `json:"package"`
	Actions map[string][]struct {
		Command string   `json:"command"`
		Args    []string `json:"args"`
	} `json:"actions"`
}

func LoadConfig(path string) (Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	config, err := ParseConfig(bytes)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func ParseConfig(configFile []byte) (Config, error) {
	var config Config
	err := json.Unmarshal(configFile, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func match(patterns []string, path string) bool {
	for _, pattern := range patterns {
		if match, _ := filepath.Match(pattern, path); match {
			return true
		}
	}
	return false
}

func (c Config) Matches(path string) bool {
	return match(c.Match, path) && !match(c.Ignore, path)
}

func (c Config) isPackageDir(dir string) bool {
	for _, filename := range c.Package {
		packageFile := filepath.Join(dir, filename)
		if fileExists(packageFile) {
			return true
		}
	}
	return false
}

func (c Config) FindPackage(path string) string {
	dir := filepath.Dir(path)
	if dir == "." || c.isPackageDir(dir) {
		return dir
	}
	return c.FindPackage(dir)
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
