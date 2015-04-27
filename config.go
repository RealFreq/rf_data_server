package main

import (
	"github.com/spf13/viper"
	"log"
)

var graphiteConfig *viper.Viper
var serverConfig *viper.Viper

func init() {
	graphiteConfig = viper.New()
	serverConfig = viper.New()
	graphiteConfig.AddConfigPath("config")
	serverConfig.AddConfigPath("config")
}

func GraphiteConfig() (string, int) {
	graphiteConfig.SetDefault("graphite", map[string]string{"host": "127.0.0.1", "port": "2003"})

	graphiteConfig.SetConfigName("logger")

	err := graphiteConfig.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading in config file: %s\n", err)
	}

	config := graphiteConfig.GetStringMap("graphite")

	host := config["host"].(string)
	port := config["port"].(int)

	return host, port
}

func ServerConfig() (string, int) {
	serverConfig.SetDefault("server", map[string]string{"host": "0.0.0.0", "port": "10000"})

	serverConfig.SetConfigName("server")

	err := serverConfig.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading in config file: %s\n", err)
	}

	config := serverConfig.GetStringMap("server")

	host := config["host"].(string)
	port := config["port"].(int)

	return host, port
}
