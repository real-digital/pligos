package helm

import (
	yaml "gopkg.in/yaml.v2"
	"realcloud.tech/pligos/pkg/pligos"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type Creator struct {
}

func (c *Creator) Save(config pligos.CreateConfig, dest string) error {
	chrt, err := c.Create(config)
	if err != nil {
		return err
	}

	return chartutil.SaveDir(chrt, dest)
}

func (c *Creator) Create(config pligos.CreateConfig) (*chart.Chart, error) {
	dependencies := []*chart.Chart{}
	for _, e := range config.Dependencies {
		dependency, err := c.Create(e)
		if err != nil {
			return nil, err
		}

		dependencies = append(dependencies, dependency)
	}

	for _, e := range config.ChartDependencies {
		dependency, err := chartutil.LoadFile(e)
		if err != nil {
			return nil, err
		}
		dependencies = append(dependencies, dependency)
	}

	base, err := chartutil.Load(config.FlavorPath)
	if err != nil {
		return nil, err
	}

	base.Metadata = &chart.Metadata{Name: config.Name, Description: config.Description, Version: config.Version}

	var updatedTemplates []*chart.Template
	for _, template := range base.Templates {
		newData := chartutil.Transform(string(template.Data), "<CHARTNAME>", base.Metadata.Name)
		updatedTemplates = append(updatedTemplates, &chart.Template{Name: template.Name, Data: newData})
	}
	base.Templates = updatedTemplates

	values, err := config.Compiler.Compile()
	if err != nil {
		return nil, err
	}

	valuesYaml, err := yaml.Marshal(values)
	if err != nil {
		return nil, err
	}

	base.Values = &chart.Config{Raw: string(valuesYaml)}
	base.Dependencies = dependencies
	return base, nil
}
