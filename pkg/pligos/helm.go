package pligos

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"

	"github.com/otiai10/copy"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type Helm struct {
	starterDir string
	outputDir  string
	starter    string
	values     map[string]interface{}

	SecretProvider SecretProvider
}

func NewHelm(sp SecretProvider, flavorDir, outputDir string, values map[string]interface{}) (*Helm, error) {
	cmd := exec.Command("helm", "home")
	helmHome, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	starterDir, err := ioutil.TempDir(filepath.Join(strings.TrimSpace(string(helmHome)), "starters"), "pligos")
	if err != nil {
		return nil, err
	}

	fmt.Println(flavorDir, starterDir)
	if err := copy.Copy(flavorDir, starterDir); err != nil {
		return nil, err
	}

	return &Helm{
		SecretProvider: sp,

		starterDir: starterDir,
		starter:    filepath.Base(starterDir),
		outputDir:  outputDir,
		values:     values,
	}, nil
}

func (h *Helm) Clean() error {
	return os.RemoveAll(h.starterDir)
}

type HelmDescription struct {
	Description string
	Version     string
}

func (h *Helm) writeValues() error {
	path := filepath.Join(h.outputDir, "values.yaml")

	buf, err := yaml.Marshal(h.values)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, buf, os.FileMode(0644))
}

func (h *Helm) Create(desc HelmDescription) error {
	cmd := exec.Command("helm", "create", "--starter", h.starter, h.outputDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	info, err := os.Stat(filepath.Join(h.outputDir, "Chart.yaml"))
	if err != nil {
		return err
	}

	if err := h.writeValues(); err != nil {
		return err
	}

	if err := h.SecretProvider.writeSecrets(); err != nil {
		return err
	}

	chartYaml, err := ioutil.ReadFile(filepath.Join(h.outputDir, "Chart.yaml"))
	if err != nil {
		return err
	}

	var md chart.Metadata
	if err := yaml.Unmarshal(chartYaml, &md); err != nil {
		return err
	}

	if desc.Description != "" {
		md.Description = desc.Description
	}
	md.Version = desc.Version

	chartYaml, err = yaml.Marshal(md)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(h.outputDir, "Chart.yaml"), chartYaml, info.Mode())
}

func (h *Helm) InstallSecrets(secrets []string) error {
	return h.copyConfiguration(secrets, "secrets")
}

func (h *Helm) InstallConfigs(configs []string) error {
	return h.copyConfiguration(configs, "config")
}

func (h *Helm) copyConfiguration(configs []string, dir string) error {
	if err := os.MkdirAll(filepath.Join(h.outputDir, dir), os.FileMode(0700)); err != nil {
		return err
	}

	for _, e := range configs {
		fileName := filepath.Base(e)

		if err := copy.Copy(e, filepath.Join(h.outputDir, dir, fileName)); err != nil {
			return err
		}
	}

	return nil
}
