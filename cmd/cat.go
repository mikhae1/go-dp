package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

// catCmd represents the cat command
var catCmd = &cobra.Command{
	Use:          "cat <env|path> [target]",
	Short:        "Read files from remote environments",
	Args:         cobra.MinimumNArgs(1),
	RunE:         catRun,
	PreRunE:      initEnvContextE,
	SilenceUsage: true,
}

func catRun(cmd *cobra.Command, args []string) error {
	var catCmd = "cat "

	catPath := ""

	// real path passed, no need to parse targets
	if len(ectx.args) > 0 {
		catPath += strings.Join(ectx.args[0:], " ")
		_, err := ectx.run.Remote(catCmd + catPath)
		return err
	}

	// multi targets
	for _, env := range ectx.targets {
		catPath := ""
		if len(ectx.args) > 0 {
			// named cat
			for _, a := range ectx.args[0:] {
				for k, l := range env.Config.Remote.Cat {
					if k == a {
						catPath += l + " "
					}
				}
			}
		} else {
			// no path specified
			for _, l := range env.Config.Remote.Cat {
				catPath += l + " "
			}
		}

		catPath = strings.TrimSpace(catPath)
		if _, err := env.Remote.Start(catCmd + catPath); err != nil {
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
	rootCmd.AddCommand(catCmd)
}
