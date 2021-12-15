package config

import (
	"2021_2_MAMBa/internal/pkg/utils/log"
	"github.com/spf13/viper"
	"strings"
)

type DbConfig struct {
	User string
	Host string
	Port int
	Pass string
	Name string
}

type ConfigParams struct {
	Db          DbConfig
	ListenPort  string
	CollectPort string
	AuthPort    string
	CSRF        bool
	CORS        bool
	Secure      bool
}

func ParseMain(configPath string) ConfigParams {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("restexample")
	if configPath != "" {
		viper.SetConfigFile(configPath)
		err := viper.ReadInConfig()
		if err != nil {
			log.Warn("failed to open config")
		}
	} else {
		log.Info("Config file is not specified.")
	}
	db := DbConfig{
		User: viper.GetString("main.db.user"),
		Host: viper.GetString("main.db.host"),
		Port: viper.GetInt("main.db.port"),
		Pass: viper.GetString("main.db.pass"),
		Name: viper.GetString("main.db.name"),
	}
	return ConfigParams{
		Db:          db,
		ListenPort:  viper.GetString("main.ports.serverPort"),
		CollectPort: viper.GetString("collections.ports.collectionsPort"),
		AuthPort:    viper.GetString("auth.ports.authPort"),
		CSRF:        viper.GetBool("main.middlewares.CSRF"),
		CORS:        viper.GetBool("main.middlewares.CORS"),
	}
}

func ParseCollections(configPath string) ConfigParams {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("restexample")
	if configPath != "" {
		viper.SetConfigFile(configPath)
		err := viper.ReadInConfig()
		if err != nil {
			log.Warn("failed to open config")
		}
	} else {
		log.Info("Config file is not specified.")
	}
	db := DbConfig{
		User: viper.GetString("collections.db.user"),
		Host: viper.GetString("collections.db.host"),
		Port: viper.GetInt("collections.db.port"),
		Pass: viper.GetString("collections.db.pass"),
		Name: viper.GetString("collections.db.name"),
	}
	return ConfigParams{
		Db:          db,
		ListenPort:  viper.GetString("main.ports.serverPort"),
		CollectPort: viper.GetString("collections.ports.collectionsPort"),
	}
}

func ParseAuth(configPath string) ConfigParams {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("restexample")
	if configPath != "" {
		viper.SetConfigFile(configPath)
		err := viper.ReadInConfig()
		if err != nil {
			log.Warn("failed to open config")
		}
	} else {
		log.Info("Config file is not specified.")
	}
	db := DbConfig{
		User: viper.GetString("auth.db.user"),
		Host: viper.GetString("auth.db.host"),
		Port: viper.GetInt("auth.db.port"),
		Pass: viper.GetString("auth.db.pass"),
		Name: viper.GetString("auth.db.name"),
	}
	return ConfigParams{
		Db:          db,
		ListenPort:  viper.GetString("main.ports.serverPort"),
		AuthPort:    viper.GetString("auth.ports.authPort"),
		Secure:      viper.GetBool("auth.cookies.secure"),
	}
}
