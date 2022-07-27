package config

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"os"
)

var Config config

type config struct {
	MachineID int64 `yaml:"machineID"`
	ApiServer struct {
		Port string `yaml:"port"`
	} `yaml:"apiServer"`
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
	args := os.Args
	var cfgName string
	if len(args) > 1 {
		cfgName = args[1]
	} else {
		cfgName = ROOT + "config/config.yaml"
	}
	log.Println(cfgName)
	bytes, err := ioutil.ReadFile(cfgName)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(bytes, &Config)
	if err != nil {
		panic(err)
	}
}
