package src

import (
	"backend/src/apis"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.GET("/api/jobs", apis.QueryJobs)
	r.GET("/api/product_lines", apis.QueryProductLines)
}
