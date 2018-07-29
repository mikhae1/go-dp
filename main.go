package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v2"
)

var (
	rootEnvDir = "./envs/"
)

func getFilesNames(rootDir string) (flist []string, err error) {
	var extentions = regexp.MustCompile(".ya?ml$")

	visit := func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && extentions.MatchString(f.Name()) {
			flist = append(flist, path)
		}

		return nil
	}

	err = filepath.Walk(rootDir, visit)

	return
}

func main() {
	fileNames, _ := getFilesNames(rootEnvDir)

	for _, fname := range fileNames {
		log.Printf("..reading %s", fname)
		data, err := ioutil.ReadFile(fname)
		if err != nil {
			log.Panic(err)
		}

		if err := yaml.Unmarshal(data, &Config); err != nil {
			log.Panicf("Unmarshal: %v", err)
		}
	}

	initEnv(Config, "orchid-master-vision")
}
