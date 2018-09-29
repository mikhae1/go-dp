// TODO:
// RunEnvByEnv
// - LocalRunEnvByEnv
// - RemoteRunEnvByEnv
//
// RunTargetByTarget
// - LocalRunTargetByTarget
// - RemoteRunTargetByTarget

package cmd

import (
	"log"
	"strings"

	"github.com/mink0/exec-cmd"
	"github.com/minkolazer/gp/config"
	"github.com/minkolazer/gp/lib"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// TODO: pass context without global variable
var (
	ectx envContext
)

type envContext struct {
	ready   bool
	args    []string
	targets []config.EnvExec
	run     EnvsCmd
}

// EnvsCmd is a wrapper on EnvExec
type EnvsCmd struct {
	Envs        []config.EnvExec
	StopOnError bool
}

// EnvsRemoteRes is a result of .Remote()
type EnvsRemoteRes struct {
	EnvName    string
	TargetName string
	Result     []execmd.ClusterRes
	Error      error
}

// EnvsLocalRes is a result of .Local()
type EnvsLocalRes struct {
	EnvName    string
	TargetName string
	Result     execmd.CmdRes
	Error      error
}

// NewEnvsCmd inits EnvsCmd with defaults
func NewEnvsCmd(envs []config.EnvExec) *EnvsCmd {
	e := EnvsCmd{
		Envs: envs,
	}

	return &e
}

// Remote runs command on every env in parallel
func (e *EnvsCmd) Remote(command string) (results []EnvsRemoteRes, err error) {
	return e.runRemote(command, true)
}

// RemoteOneByOne runs command on every env in series
func (e *EnvsCmd) RemoteOneByOne(command string) (results []EnvsRemoteRes, err error) {
	return e.runRemote(command, false)
}

func (e *EnvsCmd) runRemote(command string, parallel bool) (results []EnvsRemoteRes, err error) {
	// start command on every target
	for _, env := range e.Envs {
		res := EnvsRemoteRes{
			EnvName:    env.EnvName,
			TargetName: env.TargetName,
		}

		exec := env.Remote.Run
		if parallel {
			exec = env.Remote.Start
		}

		res.Result, res.Error = exec(command)

		// save results
		results = append(results, res)
		if e.StopOnError && res.Error != nil {
			return results, res.Error
		}
	}

	// call .Wait() for every target if parallel
	if parallel {
		for _, env := range e.Envs {
			err = env.Remote.Wait()
			if e.StopOnError && err != nil {
				return
			}
		}
	}

	return
}

// initEnvContextE detects and inits environmnets [envs] and its targets [tgts]
// 	- two arrays:
//		for each env check if it configured. if not - return error.
//			for each configured env filter selected targets and filter available targets
//
// -	one array [envs]:
//		for each env check if it in configured. if not found - return error.
//
// -	one array [tgts]:
// 		for configured env check every target. if not found - return error.
func initEnvContextE(cmd *cobra.Command, args []string) (err error) {
	log.Printf("started initEnvContextE with args: %v", args)

	if len(args) == 0 {
		return
	}

	cfgEnvs := config.GetEnvs()
	argEnvs := []string{}
	argTgts := []string{}
	envArgIndex := 0
	for i, arg := range args {
		tryArg := strings.Split(arg, ",")
		if ok, _ := lib.ArrayContainsArray(tryArg, cfgEnvs); ok {
			argEnvs = tryArg
			envArgIndex = i

			// try to find targets in the next arg with best effort
			if i < len(args)-1 {
				argTgts = strings.Split(args[i+1], ",")
			}

			break
		}
	}

	// at least one environment is needed
	if len(argEnvs) == 0 {
		err = errors.Errorf("can't find environment in args, available environments: %v", cfgEnvs)
		return
	}

	// for each configured env filter available targets
	cfgTgts := []string{}
	filTgts := map[string][]string{}
	if len(argTgts) > 0 {
		for _, env := range argEnvs {
			cfgTgts, _ = config.GetTargets(env)
			filTgts[env] = []string{}
			if ok, i := lib.ArrayContainsArray(argTgts, cfgTgts); !ok {
				if len(argEnvs) == 1 && len(argTgts) > 1 {
					err = errors.Errorf(`unknown target "%s" for "%s" env, configured targets: %v`,
						argTgts[i], env, cfgTgts)
					return
				}
			}

			for _, tgt := range argTgts {
				if i := lib.ArrayContains(cfgTgts, tgt); i == -1 {
					continue
				}

				filTgts[env] = append(filTgts[env], tgt)
			}

		}
	}

	// init args
	ectx.args = append([]string(nil), args...)
	argLen := 1
	for _, tgt := range filTgts {
		if len(tgt) > 0 {
			argLen = 2
			break
		}
	}
	// pop env and target arguments
	ectx.args = append(args[:envArgIndex], ectx.args[envArgIndex+argLen:]...)

	// init envs
	for _, env := range argEnvs {
		targets, err := config.InitEnv(env, filTgts[env])
		if err != nil {
			return err
		}
		ectx.targets = append(ectx.targets, targets...)
	}
	ectx.run = *NewEnvsCmd(ectx.targets)

	ectx.ready = true

	// FIXME: not working
	// cmd.SetArgs(ecx.args)
	// cmd.ValidArgs = ecx.args
	// rootCmd.SetArgs([]string{"gtr", "flow"})

	// TODO: print help
	// cmd.SetHelpTemplate(strings.Join(envList, ","))
	// cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
	// 	fmt.Printf("hello form dfsvsdfssdfs")
	// })

	log.Printf("initialized %v ectx.targets, ectx.args: %v", len(ectx.targets), ectx.args)

	return
}
