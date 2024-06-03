package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

// register registers all routers.
func register(r *server.Hertz) {

	customizedRegister(r)
}
