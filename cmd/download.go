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
	"realcloud.tech/pligos/cmd/download"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "update-dependencies",
	Short: "download helm dependencies",
	Long:  `This will download the requirements defined by the requiremnets.yaml inside the pligos directory`,
	Run: func(cmd *cobra.Command, args []string) {
		config := download.Config{
			PligosPath: configPath,
		}

		if err := download.Download(config); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "download: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	helmCmd.AddCommand(downloadCmd)
}
