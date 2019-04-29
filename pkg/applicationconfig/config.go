package applicationconfig

import (
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
	"realcloud.tech/pligos/pkg/pathutil"
)

type PligosConfig struct {
	Path     string `filepath:"resolve"`
	Metadata Metadata
	Context  Context

	Values map[string]interface{}
}

type Metadata struct {
	Version string   `yaml:"version"`
	Types   []string `yaml:"types" filepath:"resolve"`
}

type Context struct {
	Name         string                 `yaml:"name"`
	FlavorPath   string                 `yaml:"flavor" filepath:"resolve"`
	Dependencies []Dependency           `yaml:"dependencies"`
	Spec         map[string]interface{} `yaml:"spec"`
}

type Dependency struct {
	PligosPath string `yaml:"path" filepath:"resolve"`
	Context    string `yaml:"context"`
}

func ReadPligosConfig(pligosPath string, contextName string) (PligosConfig, error) {
	configFile, err := ioutil.ReadFile(filepath.Join(pligosPath, "pligos.yaml"))
	if err != nil {
		return PligosConfig{}, err
	}

	var applicationConfig struct {
		Metadata Metadata               `yaml:"pligos"`
		Contexts map[string]Context     `yaml:"contexts"`
		Values   map[string]interface{} `yaml:"values"`
	}
	if err := yaml.Unmarshal(configFile, &applicationConfig); err != nil {
		return PligosConfig{}, err
	}

	res := PligosConfig{
		Metadata: applicationConfig.Metadata,
		Context:  applicationConfig.Contexts[contextName],
		Values:   applicationConfig.Values,
	}

	pathutil.Resolve(&res, pligosPath)
	return res, nil
}
