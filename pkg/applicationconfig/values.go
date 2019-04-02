package applicationconfig

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
	"realcloud.tech/pligos/pkg/maputil"
)

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{}, len(a)+len(b))

	for k, v := range a {
		res[k] = v
	}

	for k, v := range b {
		res[k] = v
	}

	return res
}

func createSchema(schemaPath string, types map[string]interface{}) (map[string]interface{}, error) {
	buf, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}

	var schema map[string]interface{}
	if err := yaml.Unmarshal(buf, &schema); err != nil {
		return nil, err
	}

	return (&maputil.Normalizer{}).Normalize(mergeMaps(types, schema)), nil
}

func openTypes(types []string) (map[string]interface{}, error) {
	res := make(map[string]interface{})

	for _, e := range types {
		buf, err := ioutil.ReadFile(e)
		if err != nil {
			return nil, err
		}
		var yml map[string]interface{}
		err = yaml.Unmarshal(buf, &yml)
		if err != nil {
			return nil, err
		}
		res = mergeMaps(res, yml)
	}

	return (&maputil.Normalizer{}).Normalize(res), nil
}

func openValues(valuesPath string) (map[string]interface{}, error) {
	buf, err := ioutil.ReadFile(valuesPath)
	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{})
	if err := yaml.Unmarshal(buf, &res); err != nil {
		return nil, err
	}

	return (&maputil.Normalizer{}).Normalize(res), nil
}
