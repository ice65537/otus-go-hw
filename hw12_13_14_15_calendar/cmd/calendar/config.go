package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	Server  ServerConf
}

type ServerConf struct {
	Host    string
	Port    int
	Timeout int
}

type LoggerConf struct {
	Level string
	Depth int
}

type StorageConf struct {
	Type string
}

func GetConfig() Config {
	cfgFileParsed := strings.Split(configFile, "/") // configFile - main.go: global variable
	cfgName := cfgFileParsed[len(cfgFileParsed)-1]
	cfgPath := string(configFile[:len(configFile)-len(cfgName)])
	viper.SetConfigName(cfgName)
	viper.SetConfigType("toml")
	viper.AddConfigPath(cfgPath)
	viper.AddConfigPath(".")
	viper.AddConfigPath("../../configs")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("config reading fatal error: %w", err))
	}
	// DEFAULTS
	viper.SetDefault("logger.level", "ERROR")
	viper.SetDefault("logger.depth", 0)
	viper.SetDefault("storage.type", "memory")
	viper.SetDefault("http-server.host", "0.0.0.0")
	viper.SetDefault("http-server.port", "1234")
	viper.SetDefault("http-server.timeout", "1")
	//
	cfg := Config{
		LoggerConf{
			Level: viper.GetString("logger.level"),
			Depth: viper.GetInt("logger.depth"),
		},
		StorageConf{
			Type: viper.GetString("storage.type"),
		},
		ServerConf{
			Host:    viper.GetString("http-server.host"),
			Port:    viper.GetInt("http-server.port"),
			Timeout: viper.GetInt("http-server.timeout"),
		},
	}
	fmt.Println(cfg)
	return cfg
}
