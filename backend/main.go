package main

import (
	"backend/config"
	"backend/src"
	"context"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// 等待中断信号以优雅地关闭服务器
func gracefulShutdown(service *http.Server, exitChannel chan struct{}) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Kill)
	sig := <-quit
	config.Global.Logger.Info("Catch Signal", zap.Any("sig", sig))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := service.Shutdown(ctx); err != nil {
		config.Global.Logger.Error("Server Shutdown failed!", zap.Any("err", err.Error()))
	}
	config.Global.Logger.Info("Server Shutdown")
	close(exitChannel)
}

func main() {
	config.InitialGlobal()
	// 程序关闭前，释放数据库连接
	defer func() {
		if config.Global.DB != nil {
			db, _ := config.Global.DB.DB()
			_ = db.Close()
		}
	}()

	router := gin.Default()
	// 前端项目静态资源
	router.Static("/css", "./static/dist/css")
	router.Static("/fonts", "./static/dist/fonts")
	router.Static("/js", "./static/dist/js")
	router.StaticFile("/", "./static/dist/index.html")
	router.StaticFile("/favicon.ico", "./static/dist/favicon.ico")

	apiGroup := router.Group("/api")
	src.InitRouter(apiGroup)

	service := &http.Server{
		Addr:    config.Global.Config.Service.Addr,
		Handler: router,
	}

	go func() {
		if err := service.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			config.Global.Logger.Fatal("Server listen failed!", zap.Any("err", err))
		}
	}()

	exitChannel := make(chan struct{})
	go gracefulShutdown(service, exitChannel)

	<-exitChannel

	log.Println("Server Exit ...")
}
