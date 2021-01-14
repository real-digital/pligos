package helmport

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"realcloud.tech/pligos/pkg/compiler"
	"realcloud.tech/pligos/pkg/pligos"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
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

	for _, e := range c.Dependencies() {
		if _, err := chartutil.Save(e, filepath.Join(path, "charts")); err != nil {
			return err
		}
	}

	for _, e := range c.Templates {
		if err := ioutil.WriteFile(filepath.Join(path, e.Name), e.Data, 0644); err != nil {
			return err
		}
	}

	bytes, err := yaml.Marshal(c.Values)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(path, "values.yaml"), bytes, 0644)
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

	var updatedTemplates []*chart.File
	for _, template := range p.Flavor.Templates {
		newData := transform(string(template.Data), p.Flavor.Metadata.Name)
		updatedTemplates = append(updatedTemplates, &chart.File{Name: template.Name, Data: newData})
	}

	transformedDependencies := make([]*chart.Chart, 0, len(p.Dependencies))
	for _, e := range p.Dependencies {
		c, err := Transform(e.Pligos)
		if err != nil {
			return nil, err
		}

		if e.Alias != "" {
			c.Metadata.Name = e.Alias
		}

		transformedDependencies = append(transformedDependencies, c)
	}

	ch := &chart.Chart{
		Templates: updatedTemplates,
		Metadata:  p.Metadata,
		Files:     append(p.Flavor.Files, p.Chart.Files...),
		Values:    values,
	}

	ch.SetDependencies(append(p.Chart.Dependencies(), transformedDependencies...)...)
	return ch, nil
}

// transform performs a string replacement of the specified source for
// a given key with the replacement string
func transform(src, replacement string) []byte {
	return []byte(strings.ReplaceAll(src, "<CHARTNAME>", replacement))
}
