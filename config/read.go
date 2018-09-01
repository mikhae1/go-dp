package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/imdario/mergo"
	"github.com/minkolazer/gp/lib"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

var (
	// ConfigPath should be set by cobra once one init
	// then we will use singleton and one time config read
	ConfigPath string

	// config singleton
	config     map[string]Env
	configOnce sync.Once
)

// get config files
func getConfig() map[string]Env {
	configOnce.Do(func() {
		var err error
		if config, err = readConfig(ConfigPath); err != nil {
			lib.Exception(err)
		}
	})

	return config
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
					err = errors.Errorf("duplicate environment found in config: %s", cKey)
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
