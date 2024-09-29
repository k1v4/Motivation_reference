package cfg

import (
	"Motivation_reference/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	Listen struct {
		Type   string `yaml:"type" env-default:"port"`
		Port   string `yaml:"port" env-default:"8080"`
		BindIp string `yaml:"bind_ip" env-default:"127.0.1"`
	} `yaml:"listen"`
	Db struct {
		ConnString string `yaml:"conn_string"`
	} `yaml:"DB"`
}

var once sync.Once
var cfg *Config

func GetConfig() *Config {
	once.Do(func() {
		l := logger.GetLogger()
		l.Info("reading application configuration")

		cfg = &Config{}

		if err := cleanenv.ReadConfig("../config/config.yml", cfg); err != nil {
			desc, _ := cleanenv.GetDescription(cfg, nil)
			l.Fatalf("%s %v", desc, err)
		}
	})

	return cfg
}
