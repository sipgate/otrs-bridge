package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ReadConfig initializes the viper config and turns on config watching
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
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
	})
}
