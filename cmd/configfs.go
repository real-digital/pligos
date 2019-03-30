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

	"github.com/spf13/cobra"
	"realcloud.tech/pligos/cmd/configfs"
)

// helmfsCmd represents the helmfs command
var helmfsCmd = &cobra.Command{
	Use:   "configfs",
	Short: "mount pligos configuration, secrets locally",
	Long: `this allows mounting configurations, secrets foundin a
pligos configuration to the local filesystem. This way configuration,
which is typically mounted to a countainer inside Kubernetes can also be
used locally. This supports local application development, by
simulating the same configuration interface as inside the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := configfs.Config{
			MountPoint:  mountpoint,
			ConfigPath:  configPath,
			ContextName: contextName,
		}

		if err := configfs.ConfigFS(config); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "configfs: %v\n", err)
			os.Exit(1)
		}
	},
}

var mountpoint string

func init() {
	helmfsCmd.Flags().StringVarP(&mountpoint, "mountpoint", "m", "", "path to mountpoint")
	helmfsCmd.MarkFlagRequired("mountpoint")

	rootCmd.AddCommand(helmfsCmd)
}
