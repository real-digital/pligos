package pligos

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type PligosConfig struct {
	Secrets      SecretsConfig       `yaml:"secrets"`
	CodeGen      CodeGen             `yaml:"codeGen"`
	Name         string              `yaml:"name"`
	Description  string              `yaml:"description"`
	ChartVersion string              `yaml:"chartVersion"`
	FriendsDir   string              `yaml:"friendsLibrary"`
	TaskDir      string              `yaml:"taskfile"`
	Friends      map[string][]string `yaml:"friends"`
	Contexts     []Context           `yaml:"contexts"`
	Types        []string            `yaml:"types"`
}

func (p *PligosConfig) MakePathsAbsolute(configPath string) error {
	var absError error
	f := func(p string) string {
		abs, err := filepath.Abs(filepath.Join(configPath, p))
		if err != nil {
			absError = err
			return ""
		}

		return abs
	}

	p.FriendsDir = f(p.FriendsDir)
	p.TaskDir = f(p.TaskDir)

	for i := range p.Contexts {
		for c := range p.Contexts[i].Configs {
			p.Contexts[i].Configs[c] = f(p.Contexts[i].Configs[c])
		}

		for s := range p.Contexts[i].Secrets {
			p.Contexts[i].Secrets[s] = f(p.Contexts[i].Secrets[s])
		}

		p.Contexts[i].Flavor = f(p.Contexts[i].Flavor)
		p.Contexts[i].Output = f(p.Contexts[i].Output)
	}

	for i := range p.Types {
		p.Types[i] = f(p.Types[i])
	}

	p.CodeGen.Config.Path = f(p.CodeGen.Config.Path)
	p.CodeGen.Config.ChartPath = f(p.CodeGen.Config.ChartPath)

	return absError
}

type Context struct {
	Name    string   `yaml:"name"`
	Flavor  string   `yaml:"flavor"`
	Output  string   `yaml:"output"`
	Configs []string `yaml:"configs"`
	Secrets []string `yaml:"secrets"`
}

type CodeGen struct {
	Config struct {
		Path        string `yaml:"path"`
		Package     string `yaml:"package"`
		PackageRoot string `yaml:"packageRoot"`
		ChartPath   string `yaml:"chartPath"`
	} `yaml:"config"`
}

type SecretsConfig struct {
	Provider string `yaml:"provider"`
	Path     string `yaml:"path"`
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

	return (&Normalizer{}).Normalize(mergeMaps(types, schema)), nil

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

	return (&Normalizer{}).Normalize(res), nil
}
