package main

import (
	grpcCollections "2021_2_MAMBa/internal/app/collections"
)
const CollPort = "50040"
func main() {
	grpcCollections.RunServer(CollPort)
}