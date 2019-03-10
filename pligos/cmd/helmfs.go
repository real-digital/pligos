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

	"realcloud.tech/cloud-tools/pkg/pligos"
	"realcloud.tech/cloud-tools/pkg/pligos/helm"
	"realcloud.tech/cloud-tools/pkg/pligos/helmfs"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/spf13/cobra"
)

// helmfsCmd represents the helmfs command
var helmfsCmd = &cobra.Command{
	Use:   "helmfs",
	Short: "mount pligos configuration, secrets locally",
	Long: `this allows mounting configurations, secrets foundin a
pligos configuration to the local filesystem. This way configuration,
which is typically mounted to a countainer inside Kubernetes can also be
used locally. This supports local application development, by
simulating the same configuration interface as inside the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := fuse.Mount(
			mountpoint,
			fuse.FSName("pligos-helmfs"),
			fuse.LocalVolume(),
			fuse.VolumeName("pligos-helmfs"),
		)
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()

		helmCreator := func() *helm.Helm {
			h, err := newHelm()
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "init helm: %v\n", err)
				os.Exit(1)
			}

			return h
		}

		config, err := pligos.OpenPligosConfig(configPath)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "open pligos config: %v\n", err)
		}

		err = fs.Serve(c, helmfs.New(config.DeploymentConfig.Name, helmCreator))
		if err != nil {
			log.Fatal(err)
		}

		<-c.Ready
		if err := c.MountError; err != nil {
			log.Fatal(err)
		}
	},
}

var mountpoint string

func init() {
	helmfsCmd.Flags().StringVarP(&mountpoint, "mountpoint", "m", "", "path to mountpoint")
	helmfsCmd.MarkFlagRequired("mountpoint")

	rootCmd.AddCommand(helmfsCmd)
}
