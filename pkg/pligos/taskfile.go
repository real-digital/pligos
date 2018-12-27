package pligos

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func writeTasks(config PligosConfig, taskfile interface{}) error {
	taskfileYaml, err := yaml.Marshal(taskfile)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(config.TaskDir, os.FileMode(0700)); err != nil {
		return err
	}

	return ioutil.WriteFile(
		filepath.Join(config.TaskDir, fmt.Sprintf("Taskfile.%s.yml", config.Name)),
		taskfileYaml,
		os.FileMode(0644),
	)
}

type taskTask struct {
	Description string   `yaml:"desc"`
	Directory   string   `yaml:"dir"`
	Commands    []string `yaml:"cmds"`
}

func pligosTasks(config PligosConfig) map[string]taskTask {
	return map[string]taskTask{
		"infrastructure": taskTask{
			Description: fmt.Sprintf("generate %s infrastracture", config.Name),
			Directory:   fmt.Sprintf("infrastracture/%s", config.Name),
			Commands: []string{
				"pligos infrastructure",
			},
		},
		"codegen": taskTask{
			Description: fmt.Sprintf("generate %s config code", config.Name),
			Directory:   fmt.Sprintf("infrastracture/%s", config.Name),
			Commands: []string{
				"pligos codegen",
			},
		},
		"install_dev_dependencies": taskTask{
			Description: fmt.Sprintf("update %s helm dev dependencies", config.Name),
			Commands:    []string{
			// fmt.Sprintf("helmdep -path=%s", filepath.Join(filepath.Base(config.OutputDir), "dev", config.Name)),
			},
		},
	}
}

func taskFile(config PligosConfig) interface{} {
	tasks := make(map[string]interface{})

	for _, _ = range []interface{}{} { // range over context names
		contextName := "" // TODO set contxt name

		// helmDir := filepath.Join(filepath.Base(config.OutputDir), contextName, config.Name)

		task := taskTask{
			Description: fmt.Sprintf("%s deployment of %s service", contextName, config.Name), // TODO set context name
			Commands:    []string{
			// fmt.Sprintf("helmdep -path=%s", helmDir),
			// fmt.Sprintf("helm upgrade --install %s-%s %s", contextName, config.Name, helmDir),
			},
		}

		tasks[fmt.Sprintf("deploy_%s", contextName)] = task
	}

	for k, v := range pligosTasks(config) {
		tasks[k] = v
	}

	return struct {
		Version string                 `yaml:"version"`
		Tasks   map[string]interface{} `yaml:"tasks"`
	}{
		Version: "2",
		Tasks:   tasks,
	}
}
