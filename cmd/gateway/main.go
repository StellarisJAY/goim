package main

import "github.com/stellarisJAY/goim/internal/gateway"

func main() {
	server := gateway.Server{}
	server.Init()
	err := server.Start()
	if err != nil {
		panic(err)
	}
	select {}
}
