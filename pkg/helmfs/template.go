package helmfs

import (
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/timeconv"
)

type configKind string

const (
	secret  configKind = "Secret"
	config  configKind = "ConfigMap"
	unknown configKind = "unknown"
)

func findConfigKind(manifest string) (configKind, error) {
	var kindProbe struct {
		Kind configKind `yaml:"kind"`
	}

	if err := yaml.Unmarshal([]byte(manifest), &kindProbe); err != nil {
		return configKind(""), err
	}

	for _, e := range []configKind{secret, config} {
		if e == kindProbe.Kind {
			return e, nil
		}
	}

	return unknown, nil
}

func findResourceName(kind configKind, manifest string) (string, error) {
	var nameProbe struct {
		Metadata struct {
			Name string `yaml:"name"`
		} `yaml:"metadata"`
	}

	if err := yaml.Unmarshal([]byte(manifest), &nameProbe); err != nil {
		return "", err
	}

	name := nameProbe.Metadata.Name
	if strings.HasPrefix(nameProbe.Metadata.Name, release) {
		name = nameProbe.Metadata.Name[len(fmt.Sprintf("%s-", release)):]

	}

	return fmt.Sprintf("%s-%s", strings.ToLower(string(kind)), name), nil
}

func template(chartPath string, releaseName string, valueFiles ...string) (map[string]string, error) {
	rawVals, err := vals(valueFiles)
	if err != nil {
		return nil, err
	}
	config := &chart.Config{Raw: string(rawVals), Values: map[string]*chart.Value{}}

	options := chartutil.ReleaseOptions{
		Name:      releaseName,
		IsInstall: false,
		IsUpgrade: false,
		Time:      timeconv.Now(),
		Namespace: "default",
	}

	caps := &chartutil.Capabilities{
		APIVersions: chartutil.DefaultVersionSet,
		KubeVersion: chartutil.DefaultKubeVersion,
	}

	c, err := chartutil.Load(chartPath)
	if err != nil {
		return nil, err
	}

	vals, err := chartutil.ToRenderValuesCaps(c, config, options, caps)
	if err != nil {
		return nil, err
	}

	renderer := engine.New()
	return renderer.Render(c, vals)
}

func vals(valueFiles []string) ([]byte, error) {
	base := map[string]interface{}{}

	for _, filePath := range valueFiles {
		currentMap := map[string]interface{}{}

		var bytes []byte
		var err error
		bytes, err = ioutil.ReadFile(filePath)
		if err != nil {
			return []byte{}, err
		}

		if err := yaml.Unmarshal(bytes, &currentMap); err != nil {
			return []byte{}, fmt.Errorf("failed to parse %s: %s", filePath, err)
		}

		base = mergeValues(base, currentMap)
	}

	return yaml.Marshal(base)
}

func mergeValues(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {

		if _, exists := dest[k]; !exists {
			dest[k] = v
			continue
		}
		nextMap, ok := v.(map[string]interface{})

		if !ok {
			dest[k] = v
			continue
		}

		destMap, isMap := dest[k].(map[string]interface{})
		if !isMap {
			dest[k] = v
			continue
		}

		dest[k] = mergeValues(destMap, nextMap)
	}
	return dest
}
