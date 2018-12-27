// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

import "github.com/spf13/cobra"

// codegenCmd represents the codegen command
var codegenCmd = &cobra.Command{
	Use:   "codegen",
	Short: "Generate infrastructure related code for your service",
	Long: `To interact with the infrastructure (for example read configs, listen
on ports) the service requires an API to configuration. This command
allows creating a configuration API based on the pligos defined
configuration of the service.`,
	Run: func(cmd *cobra.Command, args []string) {

		// config, err := readPligosConfig()
		// if err != nil {
		//	fmt.Fprintf(cmd.OutOrStderr(), "read config: %v\n", err)
		//	os.Exit(1)
		// }

		// files, err := pligos.GenerateCode(config, pligos.FindContext("dev", config.Contexts))
		// if err != nil {
		//	fmt.Fprintf(cmd.OutOrStderr(), "generate code: %v\n", err)
		//	os.Exit(1)
		// }

		// if err := os.MkdirAll(config.CodeGen.Config.Path, os.FileMode(0700)); err != nil {
		//	fmt.Fprintf(cmd.OutOrStderr(), "create code dir: %v\n", err)
		//	os.Exit(1)
		// }

		// if err := pligos.WriteCodeFiles(config.CodeGen.Config.Path, files); err != nil {
		//	fmt.Fprintf(cmd.OutOrStderr(), "write code files: %v\n", err)
		//	os.Exit(1)
		// }
	},
}

func init() {
	rootCmd.AddCommand(codegenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// codegenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// codegenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
