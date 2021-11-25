package main

import (
	grpcAuth "2021_2_MAMBa/internal/app/authorization"
)

const authPort = "50041"

func main() {
	grpcAuth.RunServer(authPort)
}
