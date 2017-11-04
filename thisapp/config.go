package thisapp

import (
	"fmt"
	"gopkg.in/gcfg.v1"
)

type Configs struct {
	Settings       ConfigSetting
	Databases      ConfigDB
	CircuitBreaker ConfigCB
}
type ConfigSetting struct {
	SelfURL     string
	SelfPort    string
	PublicDir   string
	TemplateDir string
}
type ConfigDB struct {
	Conn string
	Type string
}
type ConfigCB struct {
	Enable           bool
	SuccessThreshold int
	ErrorThreshold   int
	TimeoutSec       int
}

func readConfigFrom(filePath, fileName string) (Configs, error) {
	var config Configs
	configLocation := fmt.Sprintf("%s%s", filePath, fileName)

	err := gcfg.ReadFileInto(&config, configLocation)
	if err != nil {
		err = fmt.Errorf("Failed to read config in %s. Error: %s", configLocation, err.Error())
	}
	return config, err
}

func readConfig(env string) (Configs, error) {
	fileName := fmt.Sprintf("thisapp.%s.ini", env)
	filePath := "etc/config/thisapp/"
	if env == "development" {
		filePath = "files/etc/config/thisapp/"
	}
	return readConfigFrom(filePath, fileName)
}

func (am *AppModule) GetConfig() Configs {
	return *am.config
}
