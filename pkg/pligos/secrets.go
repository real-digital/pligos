package pligos

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type SecretProvider interface {
	writeSecrets() error
}

type fileSecrets struct {
	path, chartDir string
}

type SecretProviderMap map[string]SecretProvider

func NewSecretProviderMap(secretsConfig SecretsConfig, contexts []Context) (SecretProviderMap, error) {
	res := make(map[string]SecretProvider)
	for _, context := range contexts {
		res[context.Name] = NewSecretProvider(secretsConfig, context.Output)
	}

	return SecretProviderMap(res), nil
}

func (s SecretProviderMap) Get(context string) SecretProvider {
	return s[context]
}

func NewSecretProvider(secretsConfig SecretsConfig, chartDir string) SecretProvider {
	switch secretsConfig.Provider {
	case "file":
		return &fileSecrets{path: secretsConfig.Path, chartDir: chartDir}
	default:
		return nil
	}
}

func (f *fileSecrets) writeSecrets() error {
	infos, err := ioutil.ReadDir(f.path)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Join(f.chartDir, "secrets"), os.FileMode(0700)); err != nil {
		return err
	}

	for _, e := range infos {
		buf, err := ioutil.ReadFile(filepath.Join(f.path, e.Name()))
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(filepath.Join(f.chartDir, "secrets", e.Name()), buf, e.Mode()); err != nil {
			return err
		}
	}

	names := make([]string, 0, len(infos))
	for _, e := range infos {
		names = append(names, e.Name())
	}

	return ioutil.WriteFile(filepath.Join(f.chartDir, "secrets", ".gitignore"), []byte(strings.Join(names, "\n")), os.FileMode(0644))
}
