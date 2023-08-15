package apis

import (
	"backend/config"
	"backend/src/common/response"
	"backend/src/models"
	"github.com/gin-gonic/gin"
)

func QueryProductLines(c *gin.Context) {
	productLines := models.QueryProductLines(config.Global.DB)
	var data []map[string]interface{}
	for _, v := range productLines {
		data = append(data, map[string]interface{}{
			"value": v,
			"label": v,
		})
	}
	response.Success(c, response.GetSuccess, "查询成功", data)
}
