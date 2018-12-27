// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"realcloud.tech/tools/pkg/pligos"

	"github.com/spf13/cobra"
)

// infrastructureCmd represents the infrastructure command
var infrastructureCmd = &cobra.Command{
	Use:   "infrastructure",
	Short: "generate infrastructure configuration for service",
	Long: `This generates infrastructure configuration for your
	service. The command reads in a pligos configuration file
	creates helm charts for each of the contexts configured. `,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := readPligosConfig()
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "read config: %v\n", err)
			os.Exit(1)
		}

		friendCfgs, err := readFriendCfgs(config)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "read friends: %v\n", err)
			os.Exit(1)
		}

		friend, err := pligos.NewFriend(friendCfgs)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "create friend instance: %v\n", err)
			os.Exit(1)
		}

		types, err := pligos.OpenTypes(config.Types)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "open types: %v\n", err)
			os.Exit(1)
		}

		values, err := readPligosValues()
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "open values: %v\n", err)
			os.Exit(1)
		}
		ve, err := pligos.NewValuesEncoderMap(config.Contexts, types, values)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "create value encoders: %v\n", err)
			os.Exit(1)
		}

		sp, err := pligos.NewSecretProviderMap(config.Secrets, config.Contexts)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "create secret providers: %v\n", err)
			os.Exit(1)
		}

		infra := pligos.NewInfrastructure(config, friend, ve, sp)
		if err := infra.Write(); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "write infrastructure configuration: %v\n", err)
			os.Exit(1)
		}
	},
}

func readFriendCfgs(config *pligos.PligosConfig) ([][]byte, error) {
	infos, err := ioutil.ReadDir(config.FriendsDir)
	if err != nil {
		return nil, err
	}

	res := make([][]byte, 0, len(infos))

	for _, e := range infos {
		buf, err := ioutil.ReadFile(filepath.Join(config.FriendsDir, e.Name()))
		if err != nil {
			return nil, err
		}

		res = append(res, buf)
	}

	return res, nil
}

func init() {
	rootCmd.AddCommand(infrastructureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infrastructureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infrastructureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
