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
	"io/ioutil"
	"log"
	"os"

	"github.com/minkolazer/gp/config"
	"github.com/minkolazer/gp/lib"
	"github.com/spf13/cobra"
)

var (
	vFlag bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gp <command>",
	Short: "Server automation tool",
	Long: `
Automation tool for running repetitive tasks on servers
Copyright(c) 2018`,
	// PersistentPreRun: func(cmd *cobra.Command, args []string) {}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		lib.Exception(err)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	config.ConfigPath = "./envs"
	ce := os.Getenv("CONFIG")
	if ce != "" {
		config.ConfigPath = ce
	}
	rootCmd.PersistentFlags().StringVar(&config.ConfigPath, "config", config.ConfigPath, "path config files path")
	rootCmd.PersistentFlags().BoolVarP(&vFlag, "verbose", "v", vFlag, "verbose output")

	// runs after arguments is parsed
	cobra.OnInitialize(initLogger)
}

func initLogger() {
	log.SetOutput(new(config.Logger))
	log.SetFlags(0) // remove default formatting

	if !vFlag {
		log.SetOutput(ioutil.Discard) // disable global logger output
	}
}
