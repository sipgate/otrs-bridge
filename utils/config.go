package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func ReadConfig() {
	env := os.Getenv("APPLICATION_ENV")
	if env != "production" {
		viper.SetConfigName("config.sandbox")
	} else {
		viper.SetConfigName("config")
	}
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
	})
}
