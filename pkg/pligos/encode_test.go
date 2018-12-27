package pligos

import (
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/stretchr/testify/assert"
)

func Test_encode(t *testing.T) {
	schemaYaml := MustAsset("testdata/schema.yaml")
	pligosYaml := MustAsset("testdata/pligos.yaml")
	resultYaml := MustAsset("testdata/result.yaml")

	var schema map[string]interface{}
	if err := yaml.Unmarshal(schemaYaml, &schema); err != nil {
		t.Fatalf("unmarshal schema: %v", err)
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(pligosYaml, &config); err != nil {
		t.Fatalf("unmarshal pligos config: %v", err)
	}

	var expected map[string]interface{}
	if err := yaml.Unmarshal(resultYaml, &expected); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}

	normalizer := &Normalizer{}
	ve := &ValuesEncoder{
		config: normalizer.Normalize(config),
		schema: normalizer.Normalize(schema),
	}

	schema = normalizer.Normalize(schema)
	config = normalizer.Normalize(config)

	res, err := ve.Encode(schema["context"].(map[string]interface{}), config["contexts"].(map[string]interface{})["base"].(map[string]interface{}))
	if err != nil {
		t.Fatalf("graph encode: %v", err)
	}

	assert.Equal(t, normalizer.Normalize(expected), res)
}
