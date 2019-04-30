package applicationconfig

import (
	"path/filepath"

	"realcloud.tech/pligos/pkg/maputil"
	"realcloud.tech/pligos/pkg/pligos"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

func pligosGeneratedP(c *chart.Chart) bool {
	for _, e := range c.Metadata.Keywords {
		if e == "pligosgenerated" {
			return true
		}
	}

	return false
}

func tagPligosGenerated(c *chart.Chart) {
	if pligosGeneratedP(c) {
		return
	}

	c.Metadata.Keywords = append(c.Metadata.Keywords, "pligosgenerated")
}

func Decode(config PligosConfig) (pligos.Pligos, error) {
	dependencies := make([]pligos.Pligos, 0, len(config.Context.Dependencies))
	for _, e := range config.Context.Dependencies {
		dependencyConfig, err := ReadPligosConfig(e.PligosPath, e.Context)
		if err != nil {
			return pligos.Pligos{}, err
		}

		dependency, err := Decode(dependencyConfig)
		if err != nil {
			return pligos.Pligos{}, err
		}

		dependencies = append(dependencies, dependency)
	}

	flavor, err := chartutil.Load(config.Context.FlavorPath)
	if err != nil {
		return pligos.Pligos{}, err
	}

	types, err := openTypes(config.Metadata.Types)
	if err != nil {
		return pligos.Pligos{}, err
	}

	schema, err := createSchema(filepath.Join(config.Context.FlavorPath, "schema.yaml"), types)
	if err != nil {
		return pligos.Pligos{}, err
	}

	c, err := chartutil.Load(config.Path)
	if err != nil {
		return pligos.Pligos{}, err
	}

	filteredDependencies := []*chart.Chart{}
	for _, e := range c.GetDependencies() {
		if pligosGeneratedP(e) {
			continue
		}

		filteredDependencies = append(filteredDependencies, e)
	}

	tagPligosGenerated(c)

	return pligos.Pligos{
		Chart: &chart.Chart{
			Metadata:     c.GetMetadata(),
			Files:        c.GetFiles(),
			Dependencies: filteredDependencies,
		},

		Flavor: flavor,

		ContextSpec: (&maputil.Normalizer{}).Normalize(config.Context.Spec),
		Schema:      schema["context"].(map[string]interface{}),
		Types:       schema,
		Instances:   (&maputil.Normalizer{}).Normalize(config.Values),

		Dependencies: dependencies,
	}, nil
}
