package main

import "github.com/LysetsDal/docker-api/cmd/api"

func main() {
	server := api.NewAPIServer(":4000")
	server.Run()
}
