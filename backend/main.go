package main

import (
	"backend/config"
	"backend/src"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitialGlobal()

	r := gin.Default()
	src.InitRouter(r)

	err := r.Run(config.Global.Config.Service.Addr)
	if err != nil {
		panic(err.Error())
	}
}
