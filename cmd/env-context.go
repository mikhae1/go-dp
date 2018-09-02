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

// EnvsCmd is a wrapper on ClusterSSHCmd
type EnvsCmd struct {
	Cmds        []execmd.ClusterSSHCmd
	Envs        []config.EnvExec
	StartedCmds []EnvsRes
	StopOnError bool
}

type EnvsRes struct {
	EnvName    string
	TargetName string
	Result     execmd.ClusterRes
	Error      error
}

// NewEnvsCmd inits EnvsCmd with defaults
func NewEnvsCmd(envs []config.EnvExec) *EnvsCmd {
	e := EnvsCmd{}
	e.Envs = envs
	return &e
}

// // Loop through envs and runs
// func (e *EnvsCmd) Run(command string, parallel bool) ([]EnvsRes, error) {
// 	// reset started on each new start
// 	e.StartedCmds = []EnvsRes{}
// 	for i, env := range e.Envs {
// 		res := EnvsRes{}

// 		exet := e.EnvsCmd[i].Start
// 		if !parallel {
// 			exec = e.EnvsCmd[i].Run
// 		}

// 		cres.Res, cres.Err = exec(command)

// 		// save results
// 		e.StartedCmds = append(e.StartedCmds, cres)
// 		if e.StopOnError && cres.Err != nil {
// 			return e.StartedCmds, cres.Err
// 		}
// 	}

// 	return e.StartedCmds, nil
// }
