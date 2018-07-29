package main

import (
	"log"
	"strings"

	"github.com/imdario/mergo"
)

var (
	Config map[string]Env
	env    Env
)

func getParents(envs map[string]Env, envName string) []string {
	var parents []string
	var walker func(envs map[string]Env, envName string, parents []string) []string

	walker = func(envs map[string]Env, envName string, parents []string) []string {
		if ArrayContains(parents, envName) != -1 {
			log.Panicf("circular parent reference:\n%s", strings.Join(append(parents, envName), " > "))
		}

		if _, ok := envs[envName]; !ok {
			log.Panicf("unknown parent: %s", envName)
		}

		parents = append(parents, envName)
		// prepend to reverse list of parents
		// parents = append([]string{envName}, parents...)

		if envs[envName].Parent != "" {
			envName = envs[envName].Parent
			return walker(envs, envName, parents)
		}

		// all envs should inherit Config.default
		if _, ok := Config["default"]; ok {
			parents = append(parents, "default")
		}

		return parents
	}

	return walker(envs, envName, parents)
}

func initEnv(envs map[string]Env, envName string) Env {
	log.Print("initEnv...")

	if _, ok := envs[envName]; !ok {
		log.Panicf("unknown env: %s", envName)
	}

	// TODO:
	// readyness check - if env.ready ...

	for _, e := range getParents(envs, envName) {
		if err := mergo.Merge(&env, Config[e]); err != nil {
			log.Panic(err)
		}
	}

	// merge `general` into `local` and `remote``
	if err := mergo.Merge(&env.Remote, env.General); err != nil {
		log.Panic(err)
	}

	if err := mergo.Merge(&env.Local, env.General); err != nil {
		log.Panic(err)
	}

	// TODO: templating

	// spew.Dump(env)

	return env
}
