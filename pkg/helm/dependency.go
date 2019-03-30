package helm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/downloader"
	"k8s.io/helm/pkg/getter"
	"k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type DependencyManager struct{}

// defaultKeyring returns the expanded path to the default keyring.
func (d *DependencyManager) defaultKeyring() string {
	return os.ExpandEnv("$HOME/.gnupg/pubring.gpg")
}

func (d *DependencyManager) Download(requirements, lock []byte, dest string) ([]byte, error) {
	c := &chart.Chart{
		Metadata: &chart.Metadata{Name: "download-helper"},
		Templates: []*chart.Template{
			{Name: "requirements.yaml", Data: requirements},
			{Name: "requirements.lock", Data: lock},
		},
	}

	tmpDir, err := ioutil.TempDir("", "pligos")
	if err != nil {
		return nil, err
	}

	defer os.RemoveAll(tmpDir)
	if err := chartutil.SaveDir(c, tmpDir); err != nil {
		return nil, err
	}

	man := &downloader.Manager{
		Out:        os.Stdout,
		ChartPath:  filepath.Join(tmpDir, "download-helper"),
		HelmHome:   helmpath.Home(environment.DefaultHelmHome),
		Keyring:    d.defaultKeyring(),
		SkipUpdate: false,
		Getters:    getter.All(environment.EnvSettings{Home: helmpath.Home(environment.DefaultHelmHome)}),
		Debug:      true,
	}

	if err := man.Update(); err != nil {
		return nil, err
	}

	cUpdated, err := chartutil.Load(filepath.Join(tmpDir, "download-helper"))
	if err != nil {
		return nil, err
	}

	for _, e := range cUpdated.Dependencies {
		if _, err := chartutil.Save(e, dest); err != nil {
			return nil, err
		}
	}

	for _, e := range cUpdated.Files {
		if e.GetTypeUrl() == "requirements.lock" {
			return e.GetValue(), nil
		}
	}

	return nil, fmt.Errorf("unable to locate requirements.lock")
}
