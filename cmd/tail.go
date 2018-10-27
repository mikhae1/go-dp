package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

// tailCmd represents the tail command
var tailCmd = &cobra.Command{
	Use:          "tail <env|path> [target]",
	Short:        "Stream tails from remote environments",
	Args:         cobra.MinimumNArgs(1),
	RunE:         tailRun,
	PreRunE:      initEnvContextE,
	SilenceUsage: true,
}

func tailRun(cmd *cobra.Command, args []string) error {
	var tailCmd = "tailf -n1 "

	tailPath := ""

	// real path passed, no need to parse targets
	if len(ectx.args) > 0 {
		tailPath += strings.Join(ectx.args[0:], " ")
		_, err := ectx.run.Remote(tailCmd + tailPath)
		return err
	}

	// multi targets
	for _, env := range ectx.targets {
		tailPath := ""
		if len(ectx.args) > 0 {
			// named tail
			for _, a := range ectx.args[0:] {
				for k, l := range env.Config.Remote.Log {
					if k == a {
						tailPath += l + " "
					}
				}
			}
		} else {
			// no path specified
			for _, l := range env.Config.Remote.Log {
				tailPath += l + " "
			}
		}

		tailPath = strings.TrimSpace(tailPath)
		if _, err := env.Remote.Start(tailCmd + tailPath); err != nil {
			return err
		}
	}

	// wait for results, ignoring errors
	var lastError error
	for _, env := range ectx.targets {
		if err := env.Remote.Wait(); err != nil {
			lastError = err
		}
	}

	return lastError
}

func init() {
	rootCmd.AddCommand(tailCmd)
}
