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
	"log"
	"os"

	"github.com/spf13/cobra"
	"realcloud.tech/pligos/pkg/applicationconfig"
	"realcloud.tech/pligos/pkg/helmport"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pligos",
	Short: "scalable infrastructure management",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
			return
		}

		contextName := args[0]
		pligosConfig, err := applicationconfig.ReadPligosConfig(configPath, contextName)
		if err != nil {
			log.Fatalf("read pligos configuration: %v", err)
		}

		p, err := applicationconfig.Decode(pligosConfig)
		if err != nil {
			log.Fatalf("decode pligos config: %v", err)
		}

		c, err := helmport.Transform(p)
		if err != nil {
			log.Fatalf("decode pligos config: %v", err)
		}

		if err := helmport.SwitchContext(c, pligosConfig.Path); err != nil {
			log.Fatalf("switch context: %v", err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var configPath string

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", ".", "path to pligos configuration")
}
