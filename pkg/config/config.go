package config

import (
	"flag"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

var Config config

type config struct {
	SensitiveWords []string `yaml:"sensitiveWords"`
	MachineID      int64    `yaml:"machineID"`
	TokenSecretKey string   `yaml:"tokenSecretKey"`
	ApiServer      struct {
		Port string `yaml:"port"`
	} `yaml:"apiServer"`
	MySQL struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Address  string `yaml:"address"`
		DB       string `yaml:"DB"`
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
	Nsq struct {
		LookupAddresses []string `yaml:"lookupAddresses"`
		NsqdAddress     string   `yaml:"nsqdAddress"`
	} `yaml:"nsq"`
}

const ROOT = "./"

func init() {
	var cfgName string
	flag.StringVar(&cfgName, "config", "config/config.yaml", "config file name")
	flag.Parse()
	bytes, err := ioutil.ReadFile(cfgName)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(bytes, &Config)
	if err != nil {
		panic(err)
	}
}
