package create

import (
	"fmt"
	"os"
	"path/filepath"

	"realcloud.tech/pligos/cmd/compile"
	"realcloud.tech/pligos/pkg/helm"
	"realcloud.tech/pligos/pkg/pligos"
)

type Config struct {
	ConfigPath  string
	ContextName string
	ChartPath   string
}

func Create(config Config) error {
	return createTree(config.ConfigPath, config.ContextName, config.ChartPath)
}

func NewHelm(configPath, contextName string) (*helm.Helm, error) {
	config, err := pligos.OpenPligosConfig(configPath)
	if err != nil {
		return nil, err
	}

	c, err := compile.NewCompiler(configPath, contextName)
	if err != nil {
		return nil, err
	}

	context, ok := pligos.FindContext(contextName, config.Contexts)
	if !ok {
		return nil, fmt.Errorf("no such context: %s", contextName)
	}
	return helm.New(context.Flavor, config.DeploymentConfig, config.ChartDependencies, context.Configs, context.Secrets, c), nil
}

func createTree(configPath, contextName, chartPath string) error {
	parent, err := NewHelm(configPath, contextName)
	if err != nil {
		return err
	}

	if err := parent.Create(chartPath); err != nil {
		return err
	}

	config, err := pligos.OpenPligosConfig(configPath)
	if err != nil {
		return err
	}

	context, ok := pligos.FindContext(contextName, config.Contexts)
	if !ok {
		return fmt.Errorf("no such context: %s", contextName)
	}

	for _, e := range context.Dependencies {
		config, err := pligos.OpenPligosConfig(e.Pligos)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Join(chartPath, "charts"), os.ModeDir); err != nil {
			return err
		}

		if err := createTree(e.Pligos, e.Context, filepath.Join(chartPath, "charts", config.DeploymentConfig.Name)); err != nil {
			return err
		}
	}

	return nil
}
