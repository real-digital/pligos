package applicationconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
	"realcloud.tech/pligos/pkg/pligos"

	"github.com/golang/protobuf/ptypes/any"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

func Decode(config PligosConfig, contextName string) (pligos.Pligos, error) {
	context, ok := findContext(contextName, config.Contexts)
	if !ok {
		return pligos.Pligos{}, fmt.Errorf("no such context: %s", contextName)
	}

	metadata, err := readMetadata(config.Path)
	if err != nil {
		return pligos.Pligos{}, err
	}

	flavor, err := chartutil.Load(context.FlavorPath)
	if err != nil {
		return pligos.Pligos{}, err
	}

	configurationFiles := make([]*any.Any, 0, len(context.Configs)+len(context.Secrets))
	for _, e := range append(context.Configs, context.Secrets...) {
		buf, err := ioutil.ReadFile(e)
		if err != nil {
			return pligos.Pligos{}, err
		}

		path := strings.TrimPrefix(e, config.Path)
		path = strings.TrimLeft(path, "/")

		configurationFiles = append(configurationFiles, &any.Any{TypeUrl: path, Value: buf})
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

	types, err := openTypes(config.Types)
	if err != nil {
		return pligos.Pligos{}, err
	}

	schema, err := createSchema(filepath.Join(context.FlavorPath, "schema.yaml"), types)
	if err != nil {
		return pligos.Pligos{}, err
	}

	instanceConfiguration, err := openValues(filepath.Join(config.Path, "values.yaml"))
	if err != nil {
		return pligos.Pligos{}, err
	}

	contexts := instanceConfiguration["contexts"].(map[string]interface{})[contextName].(map[string]interface{})

	dependencies := make([]pligos.Pligos, 0, len(context.Dependencies))
	for _, e := range context.Dependencies {
		dependencyConfig, err := ReadPligosConfig(e.PligosPath)
		if err != nil {
			return pligos.Pligos{}, err
		}

		dependency, err := Decode(dependencyConfig, e.Context)
		if err != nil {
			return pligos.Pligos{}, err
		}

		dependencies = append(dependencies, dependency)
	}

	return pligos.Pligos{
		Metadata:           metadata,
		Flavor:             flavor,
		ChartDependencies:  chartDependencies,
		ConfigurationFiles: configurationFiles,

		Contexts:  contexts,
		Schema:    schema["context"].(map[string]interface{}),
		Types:     schema,
		Instances: instanceConfiguration,

		Dependencies: dependencies,
	}, nil
}

func readMetadata(pligosPath string) (*chart.Metadata, error) {
	var metadata *chart.Metadata
	file, err := os.Open(filepath.Join(pligosPath, "Chart.yaml"))
	if err != nil {
		return nil, err
	}

	if err := yaml.NewDecoder(file).Decode(&metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}
