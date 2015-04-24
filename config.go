package main

import (
	"github.com/spf13/viper"
	"log"
)

func init() {
	viper.AddConfigPath("config")
}

func GraphiteConfig() (string, int) {
	viper.SetDefault("graphite", map[string]string{"host": "127.0.0.1", "port": "2003"})

	viper.SetConfigName("logger")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading in config file: %s\n", err)
	}

	config := viper.GetStringMap("graphite")

	host := config["host"].(string)
	port := config["port"].(int)

	return host, port
}

func ServerConfig() (string, int) {
	viper.SetDefault("server", map[string]string{"host": "0.0.0.0", "port": "10000"})

	viper.SetConfigName("server")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading in config file: %s\n", err)
	}

	config := viper.GetStringMap("server")

	host := config["host"].(string)
	port := config["port"].(int)

	return host, port
}
