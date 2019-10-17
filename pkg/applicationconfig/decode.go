package applicationconfig

import (
	"fmt"
	"github.com/mholt/archiver"
	"os"
	"path/filepath"
	"realcloud.tech/pligos/pkg/downloader"
	"realcloud.tech/pligos/pkg/pathutil"

	"realcloud.tech/pligos/pkg/maputil"
	"realcloud.tech/pligos/pkg/pligos"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

func pligosGeneratedP(c *chart.Chart) bool {
	for _, e := range c.Metadata.Keywords {
		if e == "pligosgenerated" {
			return true
		}
	}

	return false
}

func tagPligosGenerated(c *chart.Chart) {
	if pligosGeneratedP(c) {
		return
	}

	c.Metadata.Keywords = append(c.Metadata.Keywords, "pligosgenerated")
}

func Decode(config PligosConfig) (pligos.Pligos, error) {
	dependencies := make([]pligos.Pligos, 0, len(config.Context.Dependencies))
	for _, e := range config.Context.Dependencies {
		dependencyConfig, err := ReadPligosConfig(e.PligosPath, e.Context)
		if err != nil {
			return pligos.Pligos{}, err
		}

		dependency, err := Decode(dependencyConfig)
		if err != nil {
			return pligos.Pligos{}, err
		}

		dependencies = append(dependencies, dependency)
	}

	//Remove the existing flavor folder if exists
	if err := os.RemoveAll(filepath.Join(config.Path, "flavor")); err != nil {
		return pligos.Pligos{}, err
	}

	flavorDirCreated := false
	flavorsMainPath := filepath.Join(config.Path,"flavor")

	//Check if the flavor path is url or local directory
	if pathutil.IsValidUrl(config.Context.FlavorPath) {

		fmt.Println("Flavor is hosted on remote location")

		//Create flavor folder
		if err := os.Mkdir(flavorsMainPath, 0700); err != nil {
			return pligos.Pligos{}, err
		}

		flavorDirCreated = true

		//Download compressed flavor file
		downloadedFilePath,err := downloader.DownloadFile( flavorsMainPath, config.Context.FlavorPath)
		if err != nil {
			return pligos.Pligos{},err
		}

		//Unzip flavor file
		err = archiver.Unarchive(downloadedFilePath, flavorsMainPath)

		//Remove zipped file
		if err := os.Remove(downloadedFilePath); err != nil {
			return pligos.Pligos{},err
		}

		//Get the name of unzipped flavor folder
		file, _ := os.Open(flavorsMainPath)
		if err != nil {
			return pligos.Pligos{},err
		}

		defer file.Close()

		//Read only first directory within flavor folder
		flavorDir,_ := file.Readdir(1)

		//Update flavor path
		config.Context.FlavorPath = filepath.Join(flavorsMainPath,flavorDir[0].Name())
	}

	flavor, err := chartutil.Load(config.Context.FlavorPath)
	if err != nil {
		return pligos.Pligos{}, err
	}

	for i := range config.Metadata.Types {

		//Check if the type path is url or file
		if pathutil.IsValidUrl(config.Metadata.Types[i]) {

			//Check if flavor directory is already created
			if(!flavorDirCreated){

				if err := os.Mkdir(filepath.Join(config.Path, "flavor"), 0700); err != nil {
					return pligos.Pligos{}, err
				}
				flavorDirCreated = true
			}

			fmt.Println("\nType is hosted on remote location")

			downloadedFilePath, err := downloader.DownloadFile(flavorsMainPath, config.Metadata.Types[i])
			if err != nil {
				return pligos.Pligos{}, err
			}

			config.Metadata.Types[i] = downloadedFilePath
		}
	}

	types, err := openTypes(config.Metadata.Types)
	if err != nil {
		return pligos.Pligos{}, err
	}

	schema, err := createSchema(filepath.Join(config.Context.FlavorPath, "schema.yaml"), types)
	if err != nil {
		return pligos.Pligos{}, err
	}

	c, err := chartutil.Load(config.Path)
	if err != nil {
		return pligos.Pligos{}, err
	}

	filteredDependencies := []*chart.Chart{}
	for _, e := range c.GetDependencies() {
		if pligosGeneratedP(e) {
			continue
		}

		filteredDependencies = append(filteredDependencies, e)
	}

	tagPligosGenerated(c)

	return pligos.Pligos{
		Chart: &chart.Chart{
			Metadata:     c.GetMetadata(),
			Files:        c.GetFiles(),
			Dependencies: filteredDependencies,
		},

		Flavor: flavor,

		ContextSpec: (&maputil.Normalizer{}).Normalize(config.Context.Spec),
		Schema:      schema["context"].(map[string]interface{}),
		Types:       schema,
		Instances:   (&maputil.Normalizer{}).Normalize(config.Values),

		Dependencies: dependencies,
	}, nil
}
