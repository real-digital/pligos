package helmport

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	"realcloud.tech/pligos/pkg/compiler"
	"realcloud.tech/pligos/pkg/pligos"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

func SwitchContext(c *chart.Chart, path string) error {
	if err := os.RemoveAll(filepath.Join(path, "templates")); err != nil {
		return err
	}

	if err := os.RemoveAll(filepath.Join(path, "charts")); err != nil {
		return err
	}

	if err := os.Mkdir(filepath.Join(path, "templates"), 0700); err != nil {
		return err
	}

	if err := os.Mkdir(filepath.Join(path, "charts"), 0700); err != nil {
		return err
	}

	for _, e := range c.Dependencies {
		chartutil.Save(e, filepath.Join(path, "charts"))
	}

	for _, e := range c.Templates {
		if err := ioutil.WriteFile(filepath.Join(path, e.GetName()), e.GetData(), 0644); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(filepath.Join(path, "values.yaml"), []byte(c.GetValues().GetRaw()), 0644)
}

func Transform(p pligos.Pligos) (*chart.Chart, error) {
	compiler := compiler.New(
		p.ContextSpec,
		p.Schema,
		p.Types,
		p.Instances,
	)

	values, err := compiler.Compile()
	if err != nil {
		return nil, err
	}

	valuesYAML, err := yaml.Marshal(values)
	if err != nil {
		return nil, err
	}

	var updatedTemplates []*chart.Template
	for _, template := range p.Flavor.Templates {
		newData := chartutil.Transform(string(template.Data), "<CHARTNAME>", p.Flavor.Metadata.Name)
		updatedTemplates = append(updatedTemplates, &chart.Template{Name: template.Name, Data: newData})
	}

	transformedDependencies := make([]*chart.Chart, 0, len(p.Dependencies))
	for _, e := range p.Dependencies {
		c, err := Transform(e)
		if err != nil {
			return nil, err
		}

		transformedDependencies = append(transformedDependencies, c)
	}

	return &chart.Chart{
		Templates:    updatedTemplates,
		Metadata:     p.Metadata,
		Files:        append(p.Flavor.Files, p.Chart.Files...),
		Values:       &chart.Config{Raw: string(valuesYAML)},
		Dependencies: append(p.Chart.Dependencies, transformedDependencies...),
	}, nil
}
