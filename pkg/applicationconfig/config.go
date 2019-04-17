package applicationconfig

import (
	"bytes"
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
	"realcloud.tech/pligos/pkg/maputil"
	"realcloud.tech/pligos/pkg/pathutil"
)

type PligosConfig struct {
	Path     string `filepath:"resolve"`
	Metadata Metadata
	Context  Context
	Values   map[string]interface{}
}

type Metadata struct {
	Version    string   `yaml:"version"`
	Types      []string `yaml:"types"`
	FlavorPath string   `yaml:"flavor" filepath:"resolve"`
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
	pligosConfig, err := ioutil.ReadFile(filepath.Join(pligosPath, "pligos.yaml"))
	if err != nil {
		return PligosConfig{}, err
	}

	res := PligosConfig{}

	docs := bytes.Split(pligosConfig, []byte("---"))
	for _, e := range docs {
		kindProbe := struct {
			Kind string `yaml:"kind"`
		}{}

		if err := yaml.Unmarshal(e, &kindProbe); err != nil {
			return PligosConfig{}, err
		}

		switch kindProbe.Kind {
		case "pligos":
			var metadata Metadata
			if err := yaml.Unmarshal(e, &metadata); err != nil {
				return PligosConfig{}, err
			}
			res.Metadata = metadata
		case "context":
			var context Context
			if err := yaml.Unmarshal(e, &context); err != nil {
				return PligosConfig{}, err
			}
			if context.Name != contextName {
				continue
			}

			context.Spec = (&maputil.Normalizer{}).Normalize(context.Spec)

			res.Context = context
		case "values":
			var values map[string]interface{}
			if err := yaml.Unmarshal(e, &values); err != nil {
				return PligosConfig{}, err
			}

			delete(values, "kind")
			res.Values = (&maputil.Normalizer{}).Normalize(values)
		}
	}

	if res.Context.FlavorPath != "" {
		res.Metadata.FlavorPath = res.Context.FlavorPath
	}

	pathutil.Resolve(&res, pligosPath)
	return res, nil
}
