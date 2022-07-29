package config

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"os"
)

var Config config

type config struct {
	MachineID      int64  `yaml:"machineID"`
	TokenSecretKey string `yaml:"tokenSecretKey"`
	ApiServer      struct {
		Port string `yaml:"port"`
	} `yaml:"apiServer"`
	MySQL struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Address  string `yaml:"address"`
	} `yaml:"mysql"`
	Redis struct {
		Address      string `yaml:"address"`
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
		MaxRetry     int    `yaml:"maxRetry"`
		MaxIdleConns int    `yaml:"maxIdleConns"`
		IdleTimeout  int64  `yaml:"idleTimeout"`
	} `yaml:"redis"`
	Consul struct {
		Address string `yaml:"address"`
	} `yaml:"consul"`
	MongoDB struct {
		Hosts []string `yaml:"hosts"`
	} `yaml:"mongoDB"`
	Kafka struct {
		Addrs []string `yaml:"addrs"`
	} `yaml:"kafka"`
	RpcServer struct {
		Address string `yaml:"address"`
	} `yaml:"rpcServer"`
	WebsocketServer struct {
		Address string `yaml:"address"`
	} `yaml:"websocketServer"`
	Message struct {
		OfflineExpireTime int `yaml:"offlineExpireDays"`
	} `yaml:"message"`
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
