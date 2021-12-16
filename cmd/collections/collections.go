package main

import (
	grpcCollections "2021_2_MAMBa/internal/app/collections"
	"github.com/spf13/pflag"
)

func main() {
	var configPath string
	pflag.StringVarP(&configPath, "config", "c", "./cfg/cfg.yaml",
		"Config file path")
	pflag.Parse()
	grpcCollections.RunServer(configPath)
}
