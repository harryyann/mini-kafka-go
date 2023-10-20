package main

import "mini-kafka-go/internal/server"

func main() {
	s := server.DefaultSocketServer()
	s.Startup()
}
