package pligos

//go:generate go-bindata -pkg pligos templ/ testdata/

// import (
//	"bytes"
//	"encoding/base64"
//	"encoding/json"
//	"errors"
//	"fmt"
//	"go/format"
//	"html/template"
//	"io"
//	"io/ioutil"
//	"os"
//	"path/filepath"
//	"strings"

//	"gopkg.in/yaml.v2"

//	"realcloud.tech/tools/pkg/helm"

//	"github.com/ChimeraCoder/gojson"
//	"github.com/iancoleman/strcase"
// )

// func structName(fileName string) string {
//	exts := make([]string, 0)

//	for len(filepath.Ext(fileName)) > 0 {
//		exts = append([]string{filepath.Ext(fileName)[1:]}, exts...)
//		fileName = fileName[:len(fileName)-len(filepath.Ext(fileName))]
//	}

//	fileName = strings.Join(append([]string{fileName}, exts...), "_")
//	return strcase.ToCamel(fileName)
// }

// type SourceFile struct {
//	fileName string
//	content  []byte
// }

// func decodeConfig(input interface{}, isSecret bool) (string, error) {
//	if isSecret {
//		buf, err := base64.StdEncoding.DecodeString(input.(string))
//		return string(buf), err
//	}

//	return input.(string), nil
// }

// func typeConfigs(configs map[string]interface{}, pkg, configType string) ([]SourceFile, error) {
//	res := make([]SourceFile, 0)
//	for k, v := range configs {
//		content, err := decodeConfig(v, configType == "Secret")
//		if err != nil {
//			return nil, err
//		}

//		var parserFunc func(input io.Reader) (interface{}, error)
//		ct := contentType(content)

//		switch ct {
//		case "yaml":
//			parserFunc = gojson.ParseYaml
//		case "json":
//			parserFunc = gojson.ParseJson
//		default:
//			return nil, errors.New("invalid configuration type")
//		}

//		code, err := gojson.Generate(strings.NewReader(content), parserFunc, structName(k), pkg, []string{ct}, false, false)
//		if err != nil {
//			return nil, err
//		}

//		res = append(res, SourceFile{fileName: fmt.Sprintf("%s/%s.go", pkg, k), content: code})
//	}

//	return res, nil
// }

// func GenerateCode(config PligosConfig, context ContextConfiguration) ([]SourceFile, error) {
//	configs, err := templateConfigs(config, context, "configuration.yaml", "ConfigMap")
//	if err != nil {
//		return nil, err
//	}

//	secrets, err := templateConfigs(config, context, "secrets.yaml", "Secret")
//	if err != nil {
//		return nil, err
//	}

//	configTypes, err := typeConfigs(configs, "cfg", "ConfigMap")
//	if err != nil {
//		return nil, err
//	}

//	secretTypes, err := typeConfigs(secrets, "secret", "Secret")
//	if err != nil {
//		return nil, err
//	}

//	res := make([]SourceFile, 0)
//	res = append(res, configTypes...)
//	res = append(res, secretTypes...)

//	configInterface, err := renderConfigInterface(config, configs, secrets)
//	if err != nil {
//		return nil, err
//	}

//	formattedCode, err := format.Source([]byte(configInterface))
//	if err != nil {
//		return nil, err
//	}

//	return append(res, SourceFile{fileName: "config.go", content: formattedCode}), nil
// }

// func WriteCodeFiles(path string, files []SourceFile) error {
//	for _, e := range files {
//		if err := os.MkdirAll(filepath.Dir(filepath.Join(path, e.fileName)), os.FileMode(0700)); err != nil {
//			return err
//		}

//		if err := ioutil.WriteFile(filepath.Join(path, e.fileName), e.content, os.FileMode(0644)); err != nil {
//			return err
//		}
//	}

//	return nil
// }

// type configInterface struct {
//	Package             string
//	PackageRoot         string
//	ConfigInstances     []configInstance
//	SecretInstances     []configInstance
//	HasStructuredConfig bool
//	HasStructuredSecret bool
// }

// type configInstance struct {
//	Name          string
//	Type          string
//	ContentType   string
//	FileName      string
//	ChartPath     string
//	ConfigType    string
//	ContainerPath string
// }

// func hasStructuredData(configs map[string]interface{}, isSecret bool) (bool, error) {
//	for _, v := range configs {
//		content, err := decodeConfig(v, isSecret)
//		if err != nil {
//			return false, err
//		}

//		ct := contentType(content)
//		if ct != "binary" {
//			return true, nil
//		}
//	}

//	return false, nil
// }

// func renderConfigInterface(config PligosConfig, configs map[string]interface{}, secrets map[string]interface{}) (string, error) {
//	hasStructuredConfig, err := hasStructuredData(configs, false)
//	if err != nil {
//		return "", err
//	}

//	hasStructuredSecret, err := hasStructuredData(secrets, true)
//	if err != nil {
//		return "", err
//	}

//	ci := configInterface{
//		Package:             config.CodeGen.Config.Package,
//		PackageRoot:         config.CodeGen.Config.PackageRoot,
//		ConfigInstances:     make([]configInstance, 0, len(configs)),
//		SecretInstances:     make([]configInstance, 0, len(secrets)),
//		HasStructuredConfig: hasStructuredConfig,
//		HasStructuredSecret: hasStructuredSecret,
//	}

//	for k, v := range configs {
//		ct := contentType(v.(string))
//		t := fmt.Sprintf("%s.%s", "cfg", structName(k))
//		if ct == "binary" {
//			t = "[]byte"
//		}

//		ci.ConfigInstances = append(ci.ConfigInstances, configInstance{
//			Name:          structName(k),
//			Type:          t,
//			FileName:      k,
//			ChartPath:     config.CodeGen.Config.ChartPath,
//			ContentType:   ct,
//			ConfigType:    "ConfigMap",
//			ContainerPath: filepath.Join("/etc", config.Name, k),
//		})
//	}

//	for k, v := range secrets {
//		content, err := decodeConfig(v, true)
//		if err != nil {
//			return "", err
//		}

//		ct := contentType(content)
//		t := fmt.Sprintf("%s.%s", "secret", structName(k))
//		if ct == "binary" {
//			t = "[]byte"
//		}

//		ci.SecretInstances = append(ci.SecretInstances, configInstance{
//			Name:          structName(k),
//			Type:          t,
//			FileName:      k,
//			ChartPath:     config.CodeGen.Config.ChartPath,
//			ContentType:   ct,
//			ConfigType:    "Secret",
//			ContainerPath: filepath.Join("/etc", "secrets", k),
//		})
//	}

//	templ := MustAsset("templ/config.templ")

//	t := template.Must(template.New("configuration"), nil)
//	t, err = t.Parse(string(templ))
//	if err != nil {
//		return "", err
//	}

//	buf := bytes.NewBuffer(nil)
//	if err := t.Execute(buf, ci); err != nil {
//		return "", err
//	}

//	return buf.String(), nil
// }

// func findManifest(file, configType string) (string, error) {
//	resources := strings.Split(file, "---")
//	var probe struct {
//		Kind string `yaml:"kind"`
//		Type string `yaml:"type"`
//	}

//	for _, e := range resources {
//		if err := yaml.Unmarshal([]byte(e), &probe); err != nil {
//			return "", err
//		}

//		switch configType {
//		case "ConfigMap":
//			if probe.Kind == "ConfigMap" {
//				return e, nil
//			}

//		case "Secret":
//			if probe.Kind == "Secret" && probe.Type == "Opaque" {
//				return e, nil
//			}

//		}
//	}

//	return "", fmt.Errorf("manifest not found")
// }

// func templateConfigs(config PligosConfig, context ContextConfiguration, templateFile string, configType string) (map[string]interface{}, error) {
//	chartPath := filepath.Join(config.OutputDir, context.Name, config.Name)
//	templates, err := helm.Template(chartPath, context.Name, filepath.Join(chartPath, "values.yaml"))
//	if err != nil {
//		return nil, err
//	}

//	configurationTemplate := templates[filepath.Join(config.Name, "templates", templateFile)]

//	configurationManifest, err := findManifest(configurationTemplate, configType)
//	if err != nil {
//		return nil, err
//	}

//	manifest := struct {
//		Data map[string]interface{} `yaml:"data"`
//	}{}

//	if err := yaml.Unmarshal([]byte(configurationManifest), &manifest); err != nil {
//		return nil, err
//	}

//	return manifest.Data, nil
// }

// func contentType(input string) string {
//	out := make(map[string]interface{})

//	if err := yaml.Unmarshal([]byte(input), &out); err == nil {
//		return "yaml"
//	}

//	if err := json.Unmarshal([]byte(input), &out); err == nil {
//		return "json"
//	}

//	return "binary"
// }
