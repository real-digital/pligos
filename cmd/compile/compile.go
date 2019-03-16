package compile

import (
	"fmt"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
	"realcloud.tech/pligos/pkg/compiler"
	"realcloud.tech/pligos/pkg/pligos"
)

type Config struct {
	ConfigPath  string
	ContextName string
}

func Compile(config Config) error {
	c, err := NewCompiler(config.ConfigPath, config.ContextName)
	if err != nil {
		return err
	}

	values, err := c.Compile()
	if err != nil {
		return err
	}

	y, err := yaml.Marshal(values)
	if err != nil {
		return err
	}

	fmt.Println(string(y))
	return nil
}

func NewCompiler(configPath, contextName string) (*compiler.Compiler, error) {
	config, err := pligos.OpenPligosConfig(configPath)
	if err != nil {
		return nil, err
	}

	context, ok := pligos.FindContext(contextName, config.Contexts)
	if !ok {
		return nil, fmt.Errorf("no such context: %s", contextName)
	}

	types, err := pligos.OpenTypes(config.Types)
	if err != nil {
		return nil, err
	}

	schema, err := pligos.CreateSchema(filepath.Join(context.Flavor, "schema.yaml"), types)
	if err != nil {
		return nil, err
	}

	instanceConfiguration, err := pligos.OpenValues(filepath.Join(config.Path, "values.yaml"))
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
