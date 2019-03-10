// Copyright © 2019 real.digital
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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"

	"realcloud.tech/pligos/pkg/compiler"
	"realcloud.tech/pligos/pkg/pligos"
)

// compileCmd represents the compile command
var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "compile values.yaml file",
	Long: `This will compile your values.yaml based on the provided schema and
pligos configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := newCompiler()
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "create compiler: %v\n", err)
			os.Exit(1)
		}

		values, err := c.Compile()
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "compile: %v\n", err)
			os.Exit(1)
		}

		y, err := yaml.Marshal(values)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "stringify: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(y))
	},
}

func newCompiler() (*compiler.Compiler, error) {
	config, err := pligos.OpenPligosConfig(configPath)
	if err != nil {
		return nil, err
	}

	context, ok := pligos.FindContext(contextName, config.Contexts)
	if !ok {
		return nil, fmt.Errorf("no such context: %s", contextName)
	}

	types, err := pligos.OpenTypes(config.Types)
	if err != nil {
		return nil, err
	}

	schema, err := pligos.CreateSchema(filepath.Join(context.Flavor, "schema.yaml"), types)
	if err != nil {
		return nil, err
	}

	instanceConfiguration, err := pligos.OpenValues(filepath.Join(config.Path, "values.yaml"))
	if err != nil {
		return nil, err
	}

	return compiler.New(
		instanceConfiguration["contexts"].(map[string]interface{})[contextName].(map[string]interface{}),
		schema["context"].(map[string]interface{}),
		schema,
		instanceConfiguration,
	), nil
}

func init() {
	rootCmd.AddCommand(compileCmd)
}
