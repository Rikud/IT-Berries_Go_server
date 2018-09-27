package main

import "IT-Berries/gameServer"

func main() {
	var server gameServer.GameServer
	server.Prepare()
	server.Start()
}
