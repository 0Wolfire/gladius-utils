package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// GetString - Wrapper around viper GetString
func GetString(key string) string {
	return viper.GetString(key)
}

// SetupConfig - Sets up, watches, and registers default config
func SetupConfig(configName string, defaults map[string]string) {
	viper.SetConfigName(configName)

	base, err := GetGladiusBase()
	if err != nil {
		viper.AddConfigPath(".") // Search only for local config
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath(base) // OS specifc
	}

	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	err = viper.ReadInConfig() // Find and read the config file
	// Should probably fix this...
	if err != nil {
		if strings.HasPrefix(err.Error(), "Config File") {
			log.Printf("Cannot find config file: %s. Using defaults", err)
		} else { // Handle errors reading the config file
			panic(fmt.Errorf("Fatal error config file: %s", err))
		}
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
}

// GetGladiusBase - Returns the base directory
func GetGladiusBase() (string, error) {
	var m string
	var err error

	var base = flag.String("b", "", "The base directory for the gladius node")
	flag.Parse()

	if os.Getenv("GLADIUSBASE") == "" {
		switch runtime.GOOS {
		case "windows":
			m = filepath.Join(os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"), ".gladius")
		case "linux":
			m = os.Getenv("HOME") + "/.config/gladius"
		case "darwin":
			m = os.Getenv("HOME") + "/.config/gladius"
		default:
			m = ""
			err = errors.New("Unknown operating system, can't find gladius base directory. Set the GLADIUSBASE environment variable, or use the flag -b <base_dir> to add it manually.")
		}
	} else if *base != "" {
		m = *base
	} else {
		m = os.Getenv("GLADIUSBASE")
	}

	return m, err
}
