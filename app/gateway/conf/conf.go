package conf

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/discovery"
	consulOffice "github.com/hashicorp/consul/api"
	"github.com/kitex-contrib/config-consul/consul"
	registryConsul "github.com/kitex-contrib/registry-consul"
	"github.com/kr/pretty"
	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v2"
)

var (
	conf               *Config
	once               sync.Once
	ConsulClient       *consul.Client
	ConsulOfficeClient *consulOffice.Client
	ConsulResolver     *discovery.Resolver
)

type Config struct {
	Env string

	Hertz    Hertz    `yaml:"hertz"`
	MySQL    MySQL    `yaml:"mysql"`
	Redis    Redis    `yaml:"redis"`
	Registry Registry `yaml:"registry"`
}

type MySQL struct {
	DSN string `yaml:"dsn"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	Username string `yaml:"username"`
	DB       int    `yaml:"db"`
}
type Registry struct {
	RegistryAddress []string `yaml:"registry_address"`
	Username        string   `yaml:"username"`
	Password        string   `yaml:"password"`
}

type Hertz struct {
	Service          string `yaml:"service"`
	Address          string `yaml:"address"`
	EnablePprof      bool   `yaml:"enable_pprof"`
	EnableGzip       bool   `yaml:"enable_gzip"`
	EnableAccessLog  bool   `yaml:"enable_access_log"`
	LogLevel         string `yaml:"log_level"`
	LogFileName      string `yaml:"log_file_name"`
	LogMaxSize       int    `yaml:"log_max_size"`
	LogMaxBackups    int    `yaml:"log_max_backups"`
	LogMaxAge        int    `yaml:"log_max_age"`
	MetricsPort      string `yaml:"metrics_port"`
	ConsulHealthAddr string `yaml:"http_consul_health_addr"`
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConfRegister)
	return conf
}

func initConf() {
	prefix := "conf"
	confFileRelPath := filepath.Join(prefix, filepath.Join(GetEnv(), "conf.yaml"))
	content, err := ioutil.ReadFile(confFileRelPath)
	if err != nil {
		panic(err)
	}

	conf = new(Config)
	err = yaml.Unmarshal(content, conf)
	if err != nil {
		hlog.Error("parse yaml error - %v", err)
		panic(err)
	}
	if err := validator.Validate(conf); err != nil {
		hlog.Error("validate config error - %v", err)
		panic(err)
	}

	conf.Env = GetEnv()

	pretty.Printf("%+v\n", conf)
}

func initConfRegister() {
	consulOfficeClient, err := consulOffice.NewClient(&consulOffice.Config{Address: "192.168.3.6:8500"})
	if err != nil {
		hlog.Fatalf("consul client init failed: %v", err)
	}
	ConsulOfficeClient = consulOfficeClient

	client, err := consul.NewClient(consul.Options{
		Addr: "192.168.3.6:8500",
	})
	if err != nil {
		hlog.Fatalf("consul client init failed: %v", err)
	}

	client.RegisterConfigCallback("gateway/test.yaml", consul.AllocateUniqueID(), func(s string, cp consul.ConfigParser) {
		// map to config
		err = yaml.Unmarshal([]byte(s), &conf)
		if err != nil {
			panic(err)
		}
		hlog.Info("config updated: ", s)
	})

	consulResolver, err := registryConsul.NewConsulResolver("192.168.3.6:8500")
	if err != nil {
		panic(err)
	}

	// global
	ConsulClient = &client
	ConsulResolver = &consulResolver
}

func GetEnv() string {
	e := os.Getenv("GO_ENV")
	if len(e) == 0 {
		return "test"
	}
	return e
}

func LogLevel() hlog.Level {
	level := GetConf().Hertz.LogLevel
	switch level {
	case "trace":
		return hlog.LevelTrace
	case "debug":
		return hlog.LevelDebug
	case "info":
		return hlog.LevelInfo
	case "notice":
		return hlog.LevelNotice
	case "warn":
		return hlog.LevelWarn
	case "error":
		return hlog.LevelError
	case "fatal":
		return hlog.LevelFatal
	default:
		return hlog.LevelInfo
	}
}
