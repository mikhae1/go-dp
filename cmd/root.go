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
	"log"
	"os"

	"github.com/minkolazer/gp/config"
	"github.com/spf13/cobra"
)

var (
	env     config.EnvExec
	vFlag   bool
	envList []string // environment list for autocompletion

	configPath = "./envs/" // default path to environment files
)

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(colorInfo("[DEBUG]") + " " + string(bytes))
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gp <cmd>",
	Short: "Server automation tool",
	Long: `
Automation tool for running repetitive tasks on servers
Copyright(c) 2018`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.Printf("inside rootCmd PersistentPreRun with args: %v\n", args)

		// try to resolve args[0] and args[1] as env and target
		if len(args) == 0 {
			return
		}

		getEnv()

		for _, a := range args {
			if i := arrayContains(envList, a); i == -1 {
				break
			}

			var err error
			env, err = config.InitEnv(configPath, a)
			if err != nil {
				exception(err)
			}
		}

		// FIXME: remove args
		// args = []string{}
		// cmd.SetArgs(append([]string{"name"}, args...))

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exception(err)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	ce := os.Getenv("CONFIG")
	if ce != "" {
		configPath = ce
	}
	rootCmd.PersistentFlags().StringVar(&configPath, "config", configPath, "path config files path")
	rootCmd.PersistentFlags().BoolVarP(&vFlag, "verbose", "v", vFlag, "verbose output")

	cobra.OnInitialize(initConfig)
}

// initConfig runs after arguments is parsed
func initConfig() {
	// disable global logger output
	log.SetOutput(new(logWriter))
	log.SetFlags(0)

	if !vFlag {
		log.SetOutput(ioutil.Discard)
	}
}

// read config only once
func getEnv() *[]string {
	if len(envList) == 0 {
		var err error
		envList, err = config.GetEnvs(configPath)
		if err != nil {
			exception(err)
		}
	}

	return &envList
}
