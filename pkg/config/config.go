package config

import (
	"flag"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"testing"
)

var Config config

const (
	DevEnv     = "development"
	ProductEnv = "production"
	TestEnv    = "test"
)

type config struct {
	Environment    string   `yaml:"environment"`
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
		Database string `yaml:"database"`
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
	Etcd struct {
		Endpoints []string `yaml:"endpoints"`
	} `yaml:"etcd"`
	MessageQueue string `yaml:"messageQueue"`
	Naming       string `yaml:"naming"`
	Transfer     struct {
		SyncPushOnline bool `yaml:"syncPushOnline"`
	} `yaml:"transfer"`
	Gateway struct {
		UseJsonMsg    bool   `yaml:"useJsonMsg"`
		ConsumerGroup string `yaml:"consumerGroup"`
	} `yaml:"gateway"`
	Metrics struct {
		PromHttpAddr string `yaml:"promHttpAddr"`
	}
}

func init() {
	var cfgName string
	flag.StringVar(&cfgName, "config", "config/config.yaml", "config file name")
	testing.Init()
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
