package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
	// TODO
}

type LoggerConf struct {
	Level string
	// TODO
}

func NewConfig() Config {
	viper.SetConfigName(configFile)
	viper.SetConfigType("toml")
	viper.AddConfigPath("../configs/config.toml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("config reading fatal error: %w", err))
	}
	cfgParsed := viper.AllSettings()
	fmt.Println(cfgParsed)
	cfg := Config{}
	return cfg
}

// TODO
