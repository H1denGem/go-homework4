package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var ConfigData Config

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire string `mapstructure:"expire"`
}

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("failed to read config file: %w", err)
		log.Println("using default config and environment variables")
	} else {
		log.Println("config file loaded: config.yaml")
	}

	if err := viper.Unmarshal(&ConfigData); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}
}

func GetConfig() *Config {
	return &ConfigData
}
