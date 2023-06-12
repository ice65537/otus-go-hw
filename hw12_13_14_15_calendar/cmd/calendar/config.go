package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
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
	//
	cfg := Config{
		LoggerConf{
			Level: viper.GetString("logger.level"),
			Depth: viper.GetInt("logger.depth"),
		},
		StorageConf{
			Type: viper.GetString("storage.type"),
		},
	}
	fmt.Println(cfg)
	return cfg
}
