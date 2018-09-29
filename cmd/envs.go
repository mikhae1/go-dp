package cmd

import (
	"fmt"
	"strings"

	"github.com/minkolazer/gp/config"
	"github.com/spf13/cobra"
)

var envsCmd = &cobra.Command{
	Use:   "envs [environment]",
	Short: "Returns list of configured environments",

	Run: getEnvs,
}

func getEnvs(cmd *cobra.Command, args []string) {
	envs := config.GetEnvs()

	if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
		fmt.Print(strings.Join(envs, " "))
		return
	}

	envStr := args[0]

	for _, e := range envs {
		if envStr == e {
			targets, _ := config.GetTargets(e)
			fmt.Print(strings.Join(targets, " "))
		} else if strings.Contains(e, envStr) {
			fmt.Print(e + " ")
		}
	}
}

func init() {
	rootCmd.AddCommand(envsCmd)
}
