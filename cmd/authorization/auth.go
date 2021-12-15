package main

import (
	grpcAuth "2021_2_MAMBa/internal/app/authorization"
	"github.com/spf13/pflag"
)



func main() {
	var configPath string
	pflag.StringVarP(&configPath, "config", "c", "./cfg/cfg.yaml",
		"Config file path")
	pflag.Parse()
	grpcAuth.RunServer(configPath)
}
