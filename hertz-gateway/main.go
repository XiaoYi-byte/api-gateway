package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {

	h := server.Default(server.WithHostPorts(":8888"))

	register(h)
	h.Spin()
}
