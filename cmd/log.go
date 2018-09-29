package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:          "log <env|path> [target]",
	Short:        "Stream logs from remote environments",
	Args:         cobra.MinimumNArgs(1),
	RunE:         logRun,
	PreRunE:      initEnvContextE,
	SilenceUsage: true,
}

func logRun(cmd *cobra.Command, args []string) error {
	var tailCmd = "tailf -n1 "

	logPath := ""

	// real path passed, no need to parse targets
	if len(ectx.args) > 0 {
		logPath += strings.Join(ectx.args[0:], " ")
		_, err := ectx.run.Remote(tailCmd + logPath)
		return err
	}

	// multi targets
	for _, env := range ectx.targets {
		logPath := ""
		if len(ectx.args) > 0 {
			// named log
			for _, a := range ectx.args[0:] {
				for k, l := range env.Config.Remote.Log {
					if k == a {
						logPath += l + " "
					}
				}
			}
		} else {
			// no path specified
			for _, l := range env.Config.Remote.Log {
				logPath += l + " "
			}
		}

		logPath = strings.TrimSpace(logPath)
		if _, err := env.Remote.Start(tailCmd + logPath); err != nil {
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
	rootCmd.AddCommand(logCmd)
}
