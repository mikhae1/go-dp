package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/imdario/mergo"
	"github.com/mink0/exec-cmd"
	yaml "gopkg.in/yaml.v2"

	"github.com/pkg/errors"
)

// EnvExec is wrapper on Env with Local and Remote commands initialized
type EnvExec struct {
	Env   Env
	Local execmd.Cmd
	// only one remote command
	Remote execmd.ClusterSSHCmd
}

// InitEnv read config files, resolve parents,
// return initialized Env for envName
// TODO
// 	- templating
func InitEnv(configPath string, envName string, targetNames []string) (targets []EnvExec, err error) {
	var (
		config map[string]Env // unitialized envs from config files
	)

	log.Printf(`init env for %s:%s using: %s...`, envName, targetNames, configPath)

	if config, err = readConfig(configPath); err != nil {
		return
	}

	// check config contains requested envName
	if _, ok := config[envName]; !ok {
		err = errors.New("unknown env: " + envName)
		return
	}

	env := EnvExec{
		Env: config[envName],
	}

	// merge `General` field into `Local` and `Remote` fields
	if err = mergo.Merge(&env.Env.Remote, env.Env.Defaults); err != nil {
		return
	}
	if err = mergo.Merge(&env.Env.Local, env.Env.Defaults); err != nil {
		return
	}

	// get list of env parents
	envParents, err := getParents(config, envName)
	if err != nil {
		err = errors.Wrapf(err, "can't resolve '%s' env parents: %v", envName, envParents)
		return
	}

	if len(envParents) > 0 {
		log.Printf(`found %s parents: %v`, envName, envParents)

		// merge parents fields
		for _, e := range envParents {
			if err = mergo.Merge(&env.Env.Defaults, config[e].Defaults); err != nil {
				return
			}
			if err = mergo.Merge(&env.Env.Local, config[e].Local); err != nil {
				return
			}
			if err = mergo.Merge(&env.Env.Remote, config[e].Remote); err != nil {
				return
			}
		}

		// when parents merged, new data may appear in `Defaults` fields,
		// so merge `Defaults` fields into `Local` and `Remote`
		if err = mergo.Merge(&env.Env.Remote, env.Env.Defaults); err != nil {
			return
		}
		if err = mergo.Merge(&env.Env.Local, env.Env.Defaults); err != nil {
			return
		}
	}

	// init exec wrappers
	env.Local = *execmd.NewCmd()
	env.Remote = *execmd.NewClusterSSHCmd(env.Env.Remote.Hosts)

	if len(env.Env.Targets) == 0 || len(targetNames) == 0 {
		targets = append(targets, env)
		return
	}

	/*
		targets initialization

		best effort: target argument maybe not a target,
		throw error only if we have a list of the targets
	*/

	for _, tname := range targetNames {
		t, ok := env.Env.Targets[tname]
		if !ok {
			if len(targetNames) > 1 {
				err = errors.Errorf("unknown target %s for %s env", tname, envName)
				return
			}

			break
		}

		tEnv := Env{}
		// merge target environment
		if err = mergo.Merge(&tEnv.Defaults, t.Defaults); err != nil {
			return
		}
		if err = mergo.Merge(&tEnv.Local, tEnv.Defaults); err != nil {
			return
		}
		if err = mergo.Merge(&tEnv.Remote, tEnv.Defaults); err != nil {
			return
		}

		// merge target environment into env and override fields which are not default
		if err = mergo.Merge(&env.Env.Defaults, t.Defaults, mergo.WithOverride); err != nil {
			return
		}
		if err = mergo.Merge(&env.Env.Local, tEnv.Local, mergo.WithOverride); err != nil {
			return
		}
		if err = mergo.Merge(&env.Env.Remote, tEnv.Remote, mergo.WithOverride); err != nil {
			return
		}

		// re-init exec wrappers
		env.Remote = *execmd.NewClusterSSHCmd(env.Env.Remote.Hosts)
		for _, s := range env.Remote.SSHCmds {
			s.Cmd.PrefixStdout += tname
			s.Cmd.PrefixStderr += tname
			s.Cmd.PrefixCmd += tname
		}

		targets = append(targets, env)
	}

	return
}

// GetEnvs reads config and produce list of envs for auto completion
func GetEnvs(configPath string) (envs []string, err error) {
	config, err := readConfig(configPath)
	if err != nil {
		return nil, err
	}

	for name, env := range config {
		if !env.Hidden && name != "default" {
			envs = append(envs, name)
		}
	}

	return
}

// recursive parents search
func getParents(envs map[string]Env, envName string) (parents []string, err error) {
	var (
		walker func(envs map[string]Env, envName string) []string
	)

	walker = func(envs map[string]Env, envName string) []string {
		if arrayContains(parents, envName) != -1 {
			err = errors.New("circular parent reference:\n" + strings.Join(append(parents, envName), " > "))
			return parents
		}

		if _, ok := envs[envName]; !ok {
			err = errors.New("unknown parent: " + envName)
			return parents
		}

		parents = append(parents, envName)
		// prepend to reverse list of parents
		// parents = append([]string{envName}, parents...)

		if envs[envName].Parent != "" {
			envName = envs[envName].Parent
			return walker(envs, envName)
		}

		// all envs should inherit `default`` property if defined
		if _, ok := envs["default"]; ok {
			parents = append(parents, "default")
		}

		return parents
	}

	return walker(envs, envName)[1:], err
}

func readConfig(configPath string) (config map[string]Env, err error) {
	var yamls []string

	if configPath == "" {
		err = errors.New("no config path provided")
		return
	}

	if yamls, err = getYamlFiles(configPath); err != nil {
		return
	}

	for _, fname := range yamls {
		log.Printf("reading %s", fname)

		data := []byte{}
		if data, err = ioutil.ReadFile(fname); err != nil {
			return
		}

		yamlData := map[string]Env{}
		if err = yaml.Unmarshal(data, &yamlData); err != nil {
			err = errors.Wrap(err, "can't unmarshal config: %v")
			return
		}

		for cKey := range config {
			for yKey := range yamlData {
				if cKey == yKey {
					err = errors.New("duplicate environment found in config: " + cKey)
					return
				}
			}
		}

		if err = mergo.Merge(&config, yamlData); err != nil {
			return
		}
	}

	return
}

func getYamlFiles(path string) (flist []string, err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return
	}

	var extentions = regexp.MustCompile(".ya?ml$")

	visit := func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && extentions.MatchString(f.Name()) {
			flist = append(flist, path)
		}

		return nil
	}

	err = filepath.Walk(path, visit)

	return
}
