package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
	"gopkg.in/yaml.v2"
)

func main() {
	releases, err := releases()
	if err != nil {
		log.Fatal(err)
	}

	if err := build(releases, os.Args[1]); err != nil {
		log.Fatal(err)
	}
}

type release struct {
	os, arch, path string
}

func build(releases []release, version string) error {
	pluginYaml, err := ioutil.ReadFile("plugin.yaml")
	if err != nil {
		return err
	}

	var plugin map[string]interface{}
	if err := yaml.Unmarshal(pluginYaml, &plugin); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Join("dist", "build"), 0700); err != nil {
		return err
	}

	for _, e := range releases {
		plugin["version"] = version
		hook := fmt.Sprintf("chmod +x $HELM_PLUGIN_DIR/%s", filepath.Base(e.path))
		plugin["hooks"] = map[string]string{
			"install": hook,
			"update":  hook,
		}

		pluginYaml, err = yaml.Marshal(plugin)
		if err != nil {
			return err
		}

		pluginDist := filepath.Join("dist", "plugin.yaml")
		if err := ioutil.WriteFile(pluginDist, pluginYaml, 0644); err != nil {
			return err
		}

		d := filepath.Join("dist", "build", fmt.Sprintf("%s_%s_pligos.tar.gz", e.os, e.arch))

		if err := archiver.Archive([]string{e.path, pluginDist}, d); err != nil {
			return err
		}
	}

	return nil
}

func releases() ([]release, error) {
	base := "dist"

	res := make([]release, 0)

	osDirs, err := ioutil.ReadDir(base)
	if err != nil {
		return nil, err
	}

	for _, os := range osDirs {
		if !os.IsDir() {
			continue
		}

		archDirs, err := ioutil.ReadDir(filepath.Join(base, os.Name()))
		if err != nil {
			return nil, err
		}
		for _, arch := range archDirs {
			if !arch.IsDir() {
				continue
			}

			files, err := ioutil.ReadDir(filepath.Join(base, os.Name(), arch.Name()))
			if err != nil {
				return nil, err
			}

			if len(files) == 0 {
				continue
			}

			res = append(
				res,
				release{
					os:   os.Name(),
					arch: arch.Name(),
					path: filepath.Join(base, os.Name(), arch.Name(), files[0].Name()),
				},
			)
		}
	}

	return res, nil

}
