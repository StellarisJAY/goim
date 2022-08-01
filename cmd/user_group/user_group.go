package main

import "github.com/stellarisJAY/goim/internal/rpc/user_group"

func main() {
	user_group.Init()
	user_group.Start()
}
