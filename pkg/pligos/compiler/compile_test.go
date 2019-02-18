package compiler

//go:generate go-bindata -pkg compiler testdata/...

import (
	"testing"

	"realcloud.tech/cloud-tools/pkg/maputil"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func testCompile(schemaYaml, pligosYaml, resultYaml []byte, t *testing.T) {
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

	normalizer := &maputil.Normalizer{}
	c := &Compiler{
		config:    normalizer.Normalize(config)["contexts"].(map[string]interface{})["base"].(map[string]interface{}),
		schema:    normalizer.Normalize(schema)["context"].(map[string]interface{}),
		instances: normalizer.Normalize(config),
		types:     normalizer.Normalize(schema),
	}

	schema = normalizer.Normalize(schema)
	config = normalizer.Normalize(config)

	res, err := c.Compile()
	if err != nil {
		t.Fatalf("graph compile: %v", err)
	}

	assert.Equal(t, normalizer.Normalize(expected), res)
}

func Test_compile_a(t *testing.T) {
	schemaYaml := MustAsset("testdata/a/schema.yaml")
	pligosYaml := MustAsset("testdata/a/pligos.yaml")
	resultYaml := MustAsset("testdata/a/result.yaml")

	testCompile(schemaYaml, pligosYaml, resultYaml, t)
}

func Test_compile_b(t *testing.T) {
	schemaYaml := MustAsset("testdata/b/schema.yaml")
	pligosYaml := MustAsset("testdata/b/pligos.yaml")
	resultYaml := MustAsset("testdata/b/result.yaml")

	testCompile(schemaYaml, pligosYaml, resultYaml, t)
}
