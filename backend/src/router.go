package src

import (
	"backend/src/apis"

	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.RouterGroup) {
	router.GET("/jobs", apis.QueryJobs)
	router.GET("/product_lines", apis.QueryProductLines)
	router.GET("/get_config_parameters", apis.GetConfigParameters)
}
