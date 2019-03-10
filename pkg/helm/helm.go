package helm

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"realcloud.tech/pligos/pkg/compiler"
	"realcloud.tech/pligos/pkg/pligos"

	"github.com/otiai10/copy"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type Helm struct {
	helmStarter      *helmStarter
	deploymentConfig pligos.DeploymentConfig
	requirements     []map[string]interface{}
	configs, secrets []string

	Compiler *compiler.Compiler
}

func New(flavor string, deploymentConfig pligos.DeploymentConfig, requirements []map[string]interface{}, configs, secrets []string, compiler *compiler.Compiler) *Helm {
	return &Helm{
		helmStarter:      &helmStarter{flavorDir: flavor},
		requirements:     requirements,
		deploymentConfig: deploymentConfig,
		configs:          configs,
		secrets:          secrets,

		Compiler: compiler,
	}
}

func (h *Helm) Create(path string) error {
	return h.create(path)
}

func (h *Helm) create(path string) error {
	starter, err := h.helmStarter.init()
	if err != nil {
		return err
	}
	defer h.helmStarter.clean()

	target, err := ioutil.TempDir("", "pligos")
	if err != nil {
		return err
	}

	defer os.RemoveAll(target)

	target = filepath.Join(target, filepath.Base(path))

	cmd := exec.Command("helm", "create", "--starter", starter, target)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	if err := h.sync(target, path); err != nil {
		return err
	}

	if err := h.writeDescription(path); err != nil {
		return err
	}

	values, err := h.Compiler.Compile()
	if err != nil {
		return err
	}

	if err := h.writeValues(path, values); err != nil {
		return err
	}

	if err := h.copyConfiguration(path, append(h.configs, h.secrets...)); err != nil {
		return err
	}

	return h.writeDependencies(path)
}

func (h *Helm) writeValues(path string, values map[string]interface{}) error {
	path = filepath.Join(path, "values.yaml")

	buf, err := yaml.Marshal(values)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, buf, os.FileMode(0644))
}

func (h *Helm) copyConfiguration(path string, files []string) error {
	for _, e := range files {
		i, err := os.Stat(filepath.Dir(e))
		if err != nil {
			return err
		}

		baseDir := filepath.Base(filepath.Dir(e))
		file := filepath.Base(e)

		if err := os.MkdirAll(filepath.Join(path, baseDir), i.Mode()); err != nil {
			return err
		}

		if err := copy.Copy(e, filepath.Join(path, baseDir, file)); err != nil {
			return err
		}
	}

	return nil
}

func (h *Helm) writeDescription(path string) error {
	chartYaml, err := ioutil.ReadFile(filepath.Join(path, "Chart.yaml"))
	if err != nil {
		return err
	}

	var md chart.Metadata
	if err := yaml.Unmarshal(chartYaml, &md); err != nil {
		return err
	}

	if h.deploymentConfig.Description != "" {
		md.Description = h.deploymentConfig.Description
	}
	md.Version = h.deploymentConfig.ChartVersion

	chartYaml, err = yaml.Marshal(md)
	if err != nil {
		return err
	}

	info, err := os.Stat(filepath.Join(path, "Chart.yaml"))
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(path, "Chart.yaml"), chartYaml, info.Mode())
}

func (h *Helm) writeDependencies(path string) error {
	requirements := struct {
		Dependencies []map[string]interface{} `yaml:"dependencies"`
	}{
		Dependencies: h.requirements,
	}

	requirementsYaml, err := yaml.Marshal(requirements)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(path, "requirements.yaml"), requirementsYaml, os.FileMode(0644))
}

type helmStarter struct {
	flavorDir, starterDir string
}

func (h *helmStarter) init() (string, error) {
	t, err := ioutil.TempDir(filepath.Join(string(helmpath.Home(environment.DefaultHelmHome)), "starters"), "pligos")
	if err != nil {
		return "", err
	}

	h.starterDir = t

	if err := os.RemoveAll(t); err != nil {
		return "", err
	}

	if err := copy.Copy(h.flavorDir, h.starterDir); err != nil {
		return "", err
	}

	return filepath.Base(h.starterDir), nil
}

func (h *helmStarter) clean() error {
	return os.RemoveAll(h.starterDir)
}
