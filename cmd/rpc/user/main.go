package main

import (
	"github.com/stellarisJAY/goim/internal/rpc/user"
)

func main() {
	user.Init()
	user.Start()
}
