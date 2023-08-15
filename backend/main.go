package main

import (
	"backend/config"
	"backend/src"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitialGlobal()
	// 程序关闭前，释放数据库连接
	defer func() {
		if config.Global.DB != nil {
			db, _ := config.Global.DB.DB()
			_ = db.Close()
		}
	}()

	r := gin.Default()
	src.InitRouter(r)

	err := r.Run(config.Global.Config.Service.Addr)
	if err != nil {
		config.Global.Logger.Panic(err.Error())
		panic(err.Error())
	}
}
