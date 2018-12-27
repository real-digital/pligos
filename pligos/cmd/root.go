// Copyright © 2018 real.digital
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"realcloud.tech/tools/pkg/pligos"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pligos",
	Short: "Pligos generates infrastructure and corresponding interface to your service",
	Long: `All Pligos commands act on a single configuration (pligos.yaml) which
defines the necessary properties to host your service.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readPligosValues() (map[string]interface{}, error) {
	valuesFile, err := ioutil.ReadFile(filepath.Join(pligosPath, "values.yaml"))
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	if err := yaml.Unmarshal(valuesFile, &res); err != nil {
		return nil, err
	}

	return (&pligos.Normalizer{}).Normalize(res), nil
}

func readPligosConfig() (*pligos.PligosConfig, error) {
	cfgFile, err := ioutil.ReadFile(filepath.Join(pligosPath, "pligos.yaml"))
	if err != nil {
		return nil, err
	}

	config := new(pligos.PligosConfig)
	if err := yaml.Unmarshal(cfgFile, config); err != nil {
		return nil, err
	}
	config.MakePathsAbsolute(pligosPath)

	return config, nil
}

var pligosPath string

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pligos.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVar(&pligosPath, "path", "", "path to pligos configuration")
}
