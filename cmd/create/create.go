package create

import (
	"realcloud.tech/pligos/pkg/helm"
	"realcloud.tech/pligos/pkg/pligos"
)

type Config struct {
	ConfigPath  string
	ContextName string
	ChartPath   string
}

func Create(config Config) error {
	createConfig, err := pligos.MakeCreateConfig(config.ConfigPath, config.ContextName)
	if err != nil {
		return err
	}

	creator := &helm.Creator{}
	return creator.Save(createConfig, config.ChartPath)
}
