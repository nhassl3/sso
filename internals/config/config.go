package config

import (
	"errors"
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         uint8         `yaml:"env" env-default:"1"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

// MustLoad loading configuration of the project
// and return object in better case else
// panic and kill all program
func MustLoad() *Config {
	path, err := fetchConfigPath()
	if err != nil || path == "" {
		panic("invalid config path")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not created yet")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line or env variable
// Priority: flag > end > default
// Default: value is empty string
func fetchConfigPath() (string, error) {
	var res string

	// --config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
		if res == "" {
			return "", errors.New("CONFIG_PATH environment variable not set")
		}
		return res, nil
	}

	return res, nil
}
