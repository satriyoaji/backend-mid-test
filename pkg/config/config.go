package config

import (
	"os"

	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v3"
)

type ConfigData struct {
	Port    string `yaml:"port"`
	Env     string `yaml:"env"`
	AppCode string `yaml:"app_code"`
	Db      struct {
		Name     string `yaml:"name"`
		Host     string `yaml:"host"`
		Port     int64  `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"db"`
	RemoteTimezoneURL string `yaml:"remote_timezone_url"`
	JwtEncryptionKey  string `yaml:"jwt_encryption_key"`
}

func (c ConfigData) IsEnvProduction() bool {
	return c.Env == "production"
}

var Data *ConfigData

func Load() error {
	return LoadWithPath("./configs/config.yml")
}

func LoadWithPath(path string) error {
	if Data != nil {
		log.Info("Config already loaded")
		return nil
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(&Data); err != nil {
		return err
	}
	return nil
}
