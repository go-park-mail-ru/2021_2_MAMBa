package main

import (
	server "2021_2_MAMBa/internal/app"
	"github.com/spf13/pflag"
)

func main() {
	var configPath string
	pflag.StringVarP(&configPath, "config", "c", "./cfg/cfg.yaml",
		"Config file path")
	pflag.Parse()
	server.RunServer(configPath)
}
