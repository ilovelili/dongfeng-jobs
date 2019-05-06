package utils

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
)

var once sync.Once
var instance *Config

// GetConfig get config defined in config.json
func GetConfig() *Config {
	once.Do(func() {
		env := os.Getenv("DF_ENVIROMENT")
		if env == "" {
			env = "dev"
		}

		var config *Config
		var filepath string
		pwd, _ := os.Getwd()

		if flag.Lookup("test.v") == nil {
			// normal run
			filepath = path.Join(pwd, fmt.Sprintf("config.%s.json", strings.ToLower(env)))
		} else {
			// under go test
			filepath = path.Join(pwd, "testdata", "config.unit.test.json")
		}

		configFile, err := os.Open(filepath)
		defer configFile.Close()
		if err != nil {
			panic(err)
		}

		jsonParser := json.NewDecoder(configFile)
		err = jsonParser.Decode(&config)
		if err != nil {
			panic(err)
		}

		instance = config
	})

	return instance
}

// MySQL mysql database fields
type MySQL struct {
	Host       string `json:"host"`
	DataBase   string `json:"database"`
	User       string `json:"user"`
	Password   string `json:"password"`
	AllowDebug bool   `json:"allow_debug"`
}

// Services external services like Mysql
type Services struct {
	MySQL `json:"mysql"`
}

// ServiceMeta service meta data including service discovery specs
type ServiceMeta struct {
	Version string `json:"api_version"`
}

// Config config entry
type Config struct {
	Services    `json:"services"`
	ServiceMeta `json:"servicemeta"`
}

// GetVersion get api version
func (s *ServiceMeta) GetVersion() string {
	return s.Version
}
