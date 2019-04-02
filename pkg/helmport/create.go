package helmport

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"

	"realcloud.tech/pligos/pkg/compiler"
	"realcloud.tech/pligos/pkg/pligos"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

func Package(c *chart.Chart) ([]byte, error) {
	tmpdir, err := ioutil.TempDir("", "pligos")
	if err != nil {
		return nil, err
	}

	defer os.RemoveAll(tmpdir)

	file, err := chartutil.Save(c, tmpdir)
	if err != nil {
		panic(err)
	}

	return ioutil.ReadFile(file)
}

func Transform(p pligos.Pligos) (*chart.Chart, error) {
	compiler := compiler.New(
		p.Contexts,
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
		Files:        append(p.Flavor.Files, p.ConfigurationFiles...),
		Values:       &chart.Config{Raw: string(valuesYAML)},
		Dependencies: append(p.ChartDependencies, transformedDependencies...),
	}, nil
}
