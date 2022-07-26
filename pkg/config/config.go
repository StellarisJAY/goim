package config

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var Config config

type config struct {
	MySQL struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Address  string `yaml:"address"`
	} `yaml:"mysql"`
	Redis struct {
		Host string `yaml:"host"`
		Port uint   `yaml:"port"`
	} `yaml:"redis"`
	Consul struct {
		Address string `yaml:"address"`
	} `yaml:"consul"`
}

const ROOT = "./"

func init() {
	log.Println(getPath())
	cfgName := ROOT + "config/config.yaml"
	bytes, err := ioutil.ReadFile(cfgName)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(bytes, &Config)
	if err != nil {
		panic(err)
	}
}

func getPath() string {
	path, _ := exec.LookPath(os.Args[0])
	abs, _ := filepath.Abs(path)
	return abs
}
