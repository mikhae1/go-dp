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
	"log"
	"strings"

	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log <environment>",
	Short: "Stream logs from remote environment",
	Args:  cobra.MinimumNArgs(1),
	// ValidArgs: []string{"one", "two", "three"},
	// ValidArgs: *getEnv(),
	// Args: func(cmd *cobra.Command, args []string) error {
	// 	if len(args) < 1 {
	// 		return errors.New("one argument is required")
	// 	}
	// 	if i := arrayContains(getEnv(), args[0]); i != -1 {
	// 		return nil
	// 	}
	// 	return fmt.Errorf("unknown environment specified: %s", args[0])
	// },

	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("inside log with args: %v\n", args)

		var lpath string
		if len(args) > 1 {
			if strings.Contains(args[1], "/") {
				// real path as log targets
				lpath += strings.Join(args[1:], " ")
			} else {
				// names  log
				for _, a := range args[1:] {
					for k, l := range env.Env.Remote.Log {
						if k == a {
							lpath += l + " "
						}
					}
				}
			}
		} else {
			// no  specified
			for _, l := range env.Env.Remote.Log {
				lpath += l + " "
			}
		}

		env.Remote.Run("tail -f -n100 " + lpath)
		// env.Remote.Run("parallel --tagstring '{}:' --line-buffer tail -f {} ::: " + lpath)
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
