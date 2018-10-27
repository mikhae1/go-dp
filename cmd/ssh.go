package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// logCmd represents the log command
var sshCmd = &cobra.Command{
	Use:          "ssh <env> <host>",
	Short:        "Open ssh session to host of remote environment",
	Args:         cobra.MinimumNArgs(1),
	RunE:         sshRun,
	PreRunE:      initEnvContextE,
	SilenceUsage: true,
}

func sshRun(cmd *cobra.Command, args []string) error {
	hosts := []string{}
	paths := []string{}
	envNames := []string{}
	for _, env := range ectx.targets {
		envNames = append(envNames, env.EnvName)
		for _, cluster := range env.Remote.Cmds {
			h := cluster.Host
			if cluster.SSHCmd.User != "" {
				if !strings.Contains(h, "@") {
					h = cluster.SSHCmd.User + "@" + h
				} else {
					h = cluster.SSHCmd.User + "@" + strings.Split(h, "@")[1]
				}
			}

			hosts = append(hosts, h)
			paths = append(paths, env.Remote.Cwd)
		}
	}

	if len(ectx.args) == 0 {
		fmt.Printf("You should choose the host to ssh in %s\n", envNames)
		fmt.Println("Configured hosts:")
		for i, h := range hosts {
			fmt.Printf("%v: %s\n", i+1, h)
		}

		return fmt.Errorf("you should provide the host index or host name")
	}

	host := ""
	path := ""
	if i, err := strconv.Atoi(ectx.args[0]); err == nil {
		if i > len(hosts) {
			return fmt.Errorf("index out of range: %v for %s", i, hosts)
		}

		host = hosts[i-1]
		path = paths[i-1]
	} else {
		host = ectx.args[0]
	}

	runcmd := []string{
		host,
	}

	if path != "" {
		runcmd = append(runcmd, "-t", "cd "+path+"; exec $SHELL -l")
	}

	fmt.Printf("$ ssh %s\n", strings.Join(runcmd, " "))

	// spawn fully interactive session
	c := exec.Command("ssh", runcmd...)
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout

	return c.Run()
}

func init() {
	rootCmd.AddCommand(sshCmd)
}
