package compile

import (
	"fmt"

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
	config, err := pligos.MakePligosConfig(configPath)
	if err != nil {
		return nil, err
	}

	return pligos.MakeCompiler(config, contextName)
}
