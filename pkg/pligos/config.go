package pligos

import (
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	"realcloud.tech/pligos/pkg/maputil"
	"realcloud.tech/pligos/pkg/pathutil"
)

type PligosConfig struct {
	Version string `yaml:"version"`

	Path string `filepath:"resolve"`

	DeploymentConfig DeploymentConfig `yaml:"deployment"`

	Types []string `yaml:"types" filepath:"resolve"`

	SecretsConfig SecretsConfig `yaml:"secrets"`

	CodeGenConfig CodeGenConfig `yaml:"codeGen"`

	Contexts []Context `yaml:"contexts"`

	ChartDependencies []map[string]interface{} `yaml:"chartDependencies"`
}

func FindContext(name string, contexts []Context) (Context, bool) {
	for _, e := range contexts {
		if e.Name == name {
			return e, true
		}
	}

	return Context{}, false
}

type Dependency struct {
	Pligos  string `yaml:"pligos" filepath:"resolve"`
	Context string `yaml:"context"`
}

type DeploymentConfig struct {
	Name         string `yaml:"name"`
	Description  string `yaml:"description"`
	ChartVersion string `yaml:"chartVersion"`
}

type Context struct {
	Name         string         `yaml:"name"`
	Flavor       string         `yaml:"flavor" filepath:"resolve"`
	Configs      []string       `yaml:"configs" filepath:"resolve"`
	Secrets      []string       `yaml:"secrets" filepath:"resolve"`
	Friends      []FriendConfig `yaml:"friends"`
	Dependencies []Dependency   `yaml:"dependencies"`
}

type FriendConfig struct {
	PligosConfig string      `yaml:"path" filepath:"resolve"`
	Context      string      `yaml:"context"`
	Scope        FriendScope `yaml:"scope"`
}

type FriendScope string

const (
	Global = FriendScope("global")
	Local  = FriendScope("local")
)

type CodeGenConfig struct {
	Config struct {
		Path        string `yaml:"path" filepath:"resolve"`
		Package     string `yaml:"package"`
		PackageRoot string `yaml:"packageRoot"`
		ChartPath   string `yaml:"chartPath"`
	} `yaml:"config"`
}

type SecretsConfig struct {
	Provider string `yaml:"provider"`
	Path     string `yaml:"path" filepath:"resolve"`
}

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{}, len(a)+len(b))

	for k, v := range a {
		res[k] = v
	}

	for k, v := range b {
		res[k] = v
	}

	return res
}

func CreateSchema(schemaPath string, types map[string]interface{}) (map[string]interface{}, error) {
	buf, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}

	var schema map[string]interface{}
	if err := yaml.Unmarshal(buf, &schema); err != nil {
		return nil, err
	}

	return (&maputil.Normalizer{}).Normalize(mergeMaps(types, schema)), nil

}

func OpenPligosConfig(path string) (PligosConfig, error) {
	buf, err := ioutil.ReadFile(filepath.Join(path, "pligos.yaml"))
	if err != nil {
		return PligosConfig{}, err
	}

	var res PligosConfig
	if err := yaml.Unmarshal(buf, &res); err != nil {
		return PligosConfig{}, err
	}

	pathutil.Resolve(&res, path)
	return res, nil
}

func OpenTypes(types []string) (map[string]interface{}, error) {
	res := make(map[string]interface{})

	for _, e := range types {
		buf, err := ioutil.ReadFile(e)
		if err != nil {
			return nil, err
		}
		var yml map[string]interface{}
		err = yaml.Unmarshal(buf, &yml)
		if err != nil {
			return nil, err
		}
		res = mergeMaps(res, yml)
	}

	return (&maputil.Normalizer{}).Normalize(res), nil
}

func OpenValues(valuesPath string) (map[string]interface{}, error) {
	buf, err := ioutil.ReadFile(valuesPath)
	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{})
	if err := yaml.Unmarshal(buf, &res); err != nil {
		return nil, err
	}

	return (&maputil.Normalizer{}).Normalize(res), nil
}
