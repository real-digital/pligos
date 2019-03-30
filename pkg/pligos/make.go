package pligos

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
	"realcloud.tech/pligos/pkg/compiler"
	"realcloud.tech/pligos/pkg/pathutil"
)

func MakeCompiler(config PligosConfig, contextName string) (*compiler.Compiler, error) {
	context, ok := FindContext(contextName, config.Contexts)
	if !ok {
		return nil, fmt.Errorf("no such context: %s", contextName)
	}

	types, err := OpenTypes(config.Types)
	if err != nil {
		return nil, err
	}

	schema, err := CreateSchema(filepath.Join(context.Flavor, "schema.yaml"), types)
	if err != nil {
		return nil, err
	}

	instanceConfiguration, err := OpenValues(filepath.Join(config.Path, "values.yaml"))
	if err != nil {
		return nil, err
	}

	return compiler.New(
		instanceConfiguration["contexts"].(map[string]interface{})[contextName].(map[string]interface{}),
		schema["context"].(map[string]interface{}),
		schema,
		instanceConfiguration,
	), nil
}

func MakePligosConfig(path string) (PligosConfig, error) {
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

func MakeCreateConfig(pligosPath string, contextName string) (CreateConfig, error) {
	config, err := MakePligosConfig(pligosPath)
	if err != nil {
		return CreateConfig{}, err
	}
	context, ok := FindContext(contextName, config.Contexts)
	if !ok {
		return CreateConfig{}, fmt.Errorf("no such context: %s", contextName)
	}

	dependencies := []CreateConfig{}
	for _, e := range context.Dependencies {
		dependency, err := MakeCreateConfig(e.Pligos, e.Context)
		if err != nil {
			return CreateConfig{}, err
		}
		dependencies = append(dependencies, dependency)
	}

	chartDependencies := []string{}
	if _, err := os.Stat(filepath.Join(pligosPath, "charts")); err == nil {
		infos, err := ioutil.ReadDir(filepath.Join(pligosPath, "charts"))
		if err != nil {
			return CreateConfig{}, err
		}

		for _, e := range infos {
			chartDependencies = append(chartDependencies, filepath.Join(filepath.Join(pligosPath, "charts", e.Name())))
		}
	}

	compiler, err := MakeCompiler(config, contextName)
	if err != nil {
		return CreateConfig{}, err
	}

	return CreateConfig{
		Name:        config.DeploymentConfig.Name,
		Description: config.DeploymentConfig.Description,
		Version:     config.DeploymentConfig.ChartVersion,

		FlavorPath:         context.Flavor,
		ChartDependencies:  chartDependencies,
		ConfigurationFiles: append(context.Configs, context.Secrets...),
		Dependencies:       dependencies,
		Compiler:           compiler,
	}, nil
}
