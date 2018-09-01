// TODO:
// Run
// RunSerial

// RunEnvByEnv
// - LocalRunEnvByEnv
// - RemoteRunEnvByEnv
//
// RunTargetByTarget
// - LocalRunTargetByTarget
// - RemoteRunTargetByTarget

package cmd

import (
	"github.com/mink0/exec-cmd"
	"github.com/minkolazer/gp/config"
)

// TargetCmd is a wrapper on TargetsCmd
type EnvsCmd struct {
	Cmds []*execmd.ClusterSSHCmd

	Envs map[string][]config.EnvExec

	StartedCmds []TargetsRes
	StopOnError bool
}

type TargetsRes struct {
	env    string
	target string
	res    execmd.ClusterRes
	err    error
}

// NewTargetsCmd inits TargetsCmd with defaults
func NewTargetsCmd(targets []string) *TargetsCmd {
	t := TargetsCmd{}
	t.StopOnError = false
	t.Targets = append([]string(nil), targets...)
	for _, h := range hosts {
		t.TargetsCmd = append(t.TargetsCmd, NewSSHCmd(h))
	}
	return &t
}

// Wait wraps SSHCmd.Wait for array of hosts into t.StartedCmds struct
func (t *TargetsCmd) Wait() error {
	// TODO: list errors should be returned, or maybe hostname appended
	// now you should access `.StartedCmds` to see exact where error occurs
	var lastError error
	for i := range t.StartedCmds {
		// skip errors on Start()
		if t.StartedCmds[i].Err != nil {
			continue
		}

		t.StartedCmds[i].Err = t.TargetsCmd[i].Wait()
		if t.StartedCmds[i].Err != nil {
			if t.StopOnError {
				return t.StartedCmds[i].Err
			}

			lastError = t.StartedCmds[i].Err
		}
	}

	return lastError
}

// Run executes command in parallel: all commands starts running simultaneously at the hosts
func (t *TargetsCmd) Run(command string) (results []TargetsRes, err error) {
	if results, err = t.Start(command); err != nil {
		return
	}

	err = t.Wait()

	results = t.StartedCmds

	return
}

// RunOneByOne runs command in series: run at first host, then run at second, then...
func (t *TargetsCmd) RunOneByOne(command string) (results []TargetsRes, err error) {
	return t.start(command, false)
}

// Start runs command in parallel
func (t *TargetsCmd) Start(command string) (results []TargetsRes, err error) {
	return t.start(command, true)
}

// Loop through hosts and start
// .Start() or .Run() ssh command depending on `parallel` flag
func (t *TargetsCmd) start(command string, parallel bool) ([]TargetsRes, error) {
	// reset started on each new start
	t.StartedCmds = []TargetsRes{}
	for i, host := range t.Targets {
		cres := TargetsRes{}
		cres.Host = host

		exet := t.TargetsCmd[i].Start
		if !parallel {
			exec = t.TargetsCmd[i].Run
		}

		cres.Res, cres.Err = exec(command)

		// save results
		t.StartedCmds = append(t.StartedCmds, cres)
		if t.StopOnError && cres.Err != nil {
			return t.StartedCmds, cres.Err
		}
	}

	return t.StartedCmds, nil
}
