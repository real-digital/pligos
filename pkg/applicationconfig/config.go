package applicationconfig

import (
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
	"realcloud.tech/pligos/pkg/pathutil"
)

type PligosConfig struct {
	Path string `filepath:"resolve"`

	Version  string    `yaml:"apiVersion"`
	Types    []string  `yaml:"types" filepath:"resolve"`
	Contexts []Context `yaml:"contexts"`
}

type Context struct {
	Name         string       `yaml:"name"`
	FlavorPath   string       `yaml:"flavor" filepath:"resolve"`
	Configs      []string     `yaml:"configs" filepath:"resolve"`
	Secrets      []string     `yaml:"secrets" filepath:"resolve"`
	Dependencies []Dependency `yaml:"dependencies"`
}

type Dependency struct {
	PligosPath string `yaml:"path" filepath:"resolve"`
	Context    string `yaml:"context"`
}

func ReadPligosConfig(pligosPath string) (PligosConfig, error) {
	f, err := os.Open(filepath.Join(pligosPath, "pligos.yaml"))
	if err != nil {
		return PligosConfig{}, err
	}

	var config PligosConfig
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return PligosConfig{}, err
	}

	pathutil.Resolve(&config, pligosPath)
	return config, nil
}

func findContext(name string, contexts []Context) (Context, bool) {
	for _, e := range contexts {
		if e.Name == name {
			return e, true
		}
	}

	return Context{}, false
}
