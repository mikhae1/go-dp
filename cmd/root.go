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
	"strings"

	"github.com/minkolazer/gp/config"
	"github.com/minkolazer/gp/lib"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	ctx   envContext
	vFlag bool
	// envList []string // environment list for autocompletion

	// configPath = "./envs" // default path to environment files
	// configPath = config.ConfigPath
)

type envContext struct {
	ready   bool
	args    []string
	targets []config.EnvExec
}

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(lib.ColorInfo("[DEBUG]") + " " + string(bytes))
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gp <cmd>",
	Short: "Server automation tool",
	Long: `
Automation tool for running repetitive tasks on servers
Copyright(c) 2018`,
	PersistentPreRunE: initEnvContext,
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

	cobra.OnInitialize(initCobra)
}

// initConfig runs after arguments is parsed
func initCobra() {
	// disable global logger output
	log.SetOutput(new(logWriter))
	log.SetFlags(0)

	if !vFlag {
		log.SetOutput(ioutil.Discard)
	}
}

// // read config only once
// func getEnv() *[]string {
// 	if len(envList) == 0 {
// 		var err error
// 		envList, err = config.GetEnvs()
// 		if err != nil {
// 			lib.Exception(err)
// 		}
// 	}

// 	return &envList
// }

func initEnvContext(cmd *cobra.Command, args []string) (err error) {
	log.Printf("inside rootCmd PersistentPreRun with args: %v\n", args)

	// try to resolve args as env and targets
	if len(args) == 0 {
		return
	}

	envList := config.GetEnvs()
	// find env
	envName := ""
	ctx.args = append([]string(nil), args...) // full slice copy
	for i, arg := range args {
		if i := lib.ArrayContains(envList, arg); i == -1 {
			continue
		}

		envName = arg

		// try to find targets in the next arg
		tlist := []string{}
		if i < len(args)-1 {
			tlist = strings.Split(args[i+1], ",")
		}

		// wrong target check
		isWrongTargets := false
		wrongTarget := ""
		knownTargets, _ := config.GetTargets(envName)
		if len(tlist) > 0 {
			for _, t := range tlist {
				if lib.ArrayContains(knownTargets, t) == -1 {
					isWrongTargets = true
					wrongTarget = t
					break
				}
			}
		} else {
			isWrongTargets = true
		}

		if isWrongTargets {
			if len(tlist) > 1 {
				err = errors.Errorf(`unknown target "%s" for "%s" environment`, wrongTarget, envName)
				return
			}

			// no targets in arg+1 :(
			tlist = []string{}
		} else {
			// pop targets arg
			ctx.args = append(ctx.args[:i+1], ctx.args[i+2:]...)
		}

		if ctx.targets, err = config.InitEnv(arg, tlist); err != nil {
			return err
		}

		// pop env arg
		ctx.args = append(ctx.args[:i], ctx.args[i+1:]...)
		ctx.ready = true
		break
	}

	if !ctx.ready {
		err = errors.Errorf("unknown environment: %v", args)
		return
	}

	return nil
}
