package main

import (
	"github.com/stellarisJAY/goim/internal/rpc/auth"
	"log"
)

func main() {
	auth.Init()
	err := auth.Start()
	if err != nil {
		log.Println(err)
	}
}
