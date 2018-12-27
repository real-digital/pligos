package helm

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/timeconv"
)

func Template(chartPath string, releaseName string, valueFiles ...string) (map[string]string, error) {
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

	// User specified a values files via -f/--values
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
		// Merge with the previous map
		base = mergeValues(base, currentMap)
	}

	return yaml.Marshal(base)
}

func mergeValues(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = v
			continue
		}
		nextMap, ok := v.(map[string]interface{})
		// If it isn't another map, overwrite the value
		if !ok {
			dest[k] = v
			continue
		}
		// Edge case: If the key exists in the destination, but isn't a map
		destMap, isMap := dest[k].(map[string]interface{})
		// If the source map has a map for this key, prefer it
		if !isMap {
			dest[k] = v
			continue
		}
		// If we got to this point, it is a map in both, so merge them
		dest[k] = mergeValues(destMap, nextMap)
	}
	return dest
}
