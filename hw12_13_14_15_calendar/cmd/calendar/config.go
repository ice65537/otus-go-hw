package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	ServerHTTP  ServerConf
	ServerGRPC  ServerConf
}

type StorageConf struct {
	Type    string
	Postgre StorePostgresConf
}

type StorePostgresConf struct {
	Host     string
	Port     int
	Dbname   string
	Username string
	Password string
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

func GetConfig() Config {
	cfgFileParsed := strings.Split(configFile, "/") // configFile - main.go: global variable
	cfgName := cfgFileParsed[len(cfgFileParsed)-1]
	cfgPath := configFile[:len(configFile)-len(cfgName)]
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

	viper.SetDefault("grpc-server.host", "0.0.0.0")
	viper.SetDefault("grpc-server.port", "1235")
	viper.SetDefault("grpc-server.timeout", "1")

	viper.SetDefault("store-postgres.host", "localhost")
	viper.SetDefault("store-postgres.port", "2345")
	viper.SetDefault("store-postgres.dbname", "calendar")
	viper.SetDefault("store-postgres.username", "clndr")
	viper.SetDefault("store-postgres.password", "")
	//
	cfg := Config{
		Logger: LoggerConf{
			Level: viper.GetString("logger.level"),
			Depth: viper.GetInt("logger.depth"),
		},
		Storage: StorageConf{
			Type: viper.GetString("storage.type"),
		},
		ServerHTTP: ServerConf{
			Host:    viper.GetString("http-server.host"),
			Port:    viper.GetInt("http-server.port"),
			Timeout: viper.GetInt("http-server.timeout"),
		},
		ServerGRPC: ServerConf{
			Host:    viper.GetString("grpc-server.host"),
			Port:    viper.GetInt("grpc-server.port"),
			Timeout: viper.GetInt("grpc-server.timeout"),
		},
	}

	if cfg.Storage.Type == "postgres" {
		cfg.Storage.Postgre = StorePostgresConf{
			Host:     viper.GetString("store-postgres.host"),
			Port:     viper.GetInt("store-postgres.port"),
			Dbname:   viper.GetString("store-postgres.dbname"),
			Username: viper.GetString("store-postgres.username"),
			Password: viper.GetString("store-postgres.password"),
		}
	}
	return cfg
}
