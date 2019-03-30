package download

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"realcloud.tech/pligos/pkg/helm"
)

type Config struct {
	PligosPath string
}

func Download(config Config) error {
	_, err := os.Stat(filepath.Join(config.PligosPath, "requirements.yaml"))
	if err != nil {
		return err
	}

	requirements, err := ioutil.ReadFile(filepath.Join(config.PligosPath, "requirements.yaml"))
	if err != nil {
		return err
	}

	var lock []byte
	if _, err := os.Stat(filepath.Join(config.PligosPath, "requirements.lock")); err == nil {
		lock, err = ioutil.ReadFile(filepath.Join(config.PligosPath, "requirements.lock"))
		if err != nil {
			return err
		}
	}

	if err := os.MkdirAll(filepath.Join(config.PligosPath, "charts"), 0755); err != nil {
		return err
	}

	manager := &helm.DependencyManager{}
	updatedLock, err := manager.Download(requirements, lock, filepath.Join(config.PligosPath, "charts"))
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(config.PligosPath, "requirements.lock"), updatedLock, 0644)
}
