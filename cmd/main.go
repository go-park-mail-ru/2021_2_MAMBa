package main

import server "2021_2_MAMBa/internal/app"

const serverPort = ":8080"

func main() {
	server.RunServer(serverPort)
}
