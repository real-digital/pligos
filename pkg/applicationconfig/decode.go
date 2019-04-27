package applicationconfig

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"realcloud.tech/pligos/pkg/pligos"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

func Decode(config PligosConfig) (pligos.Pligos, error) {
	flavor, err := chartutil.Load(config.Context.FlavorPath)
	if err != nil {
		return pligos.Pligos{}, err
	}

	var chartDependencies []*chart.Chart
	if _, err := os.Stat(filepath.Join(config.Path, "charts")); err == nil {
		infos, err := ioutil.ReadDir(filepath.Join(config.Path, "charts"))
		if err != nil {
			return pligos.Pligos{}, err
		}

		for _, e := range infos {
			next, err := chartutil.Load(filepath.Join(config.Path, "charts", e.Name()))
			if err != nil {
				return pligos.Pligos{}, err
			}

			chartDependencies = append(chartDependencies, next)
		}
	}

	types, err := openTypes(config.Metadata.Types)
	if err != nil {
		return pligos.Pligos{}, err
	}

	schema, err := createSchema(filepath.Join(config.Context.FlavorPath, "schema.yaml"), types)
	if err != nil {
		return pligos.Pligos{}, err
	}

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

	c, err := chartutil.Load(config.Path)
	if err != nil {
		return pligos.Pligos{}, err
	}

	return pligos.Pligos{
		Chart: &chart.Chart{
			Metadata:     c.GetMetadata(),
			Files:        c.GetFiles(),
			Dependencies: c.GetDependencies(),
		},

		Flavor: flavor,

		ContextSpec: config.Context.Spec,
		Schema:      schema["context"].(map[string]interface{}),
		Types:       schema,
		Instances:   config.Values,

		Dependencies: dependencies,
	}, nil
}
