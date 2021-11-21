package main

import (
	server "2021_2_MAMBa/internal/app"
)

const serverPort = ":8080"
const collPort = "50040"
const authPort = "50041"

func main() {
	server.RunServer(serverPort, collPort, authPort)
}
