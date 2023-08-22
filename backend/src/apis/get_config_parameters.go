package apis

import (
	"backend/config"
	"backend/src/auth"
	"backend/src/common/response"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

func generateNonceStr() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	r.Seed(time.Now().UnixNano())
	b := make([]rune, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetConfigParameters(c *gin.Context) {
	nonceStr := generateNonceStr()
	url := c.Query("url")
	ticket, err := auth.GetTicket()
	if err != nil {
		response.Fail(c, response.GetFailed, err.Error())
		return
	}
	timestamp := time.Now().Unix() * 1000
	verifyStr := fmt.Sprintf(
		"jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		ticket, nonceStr, timestamp, url,
	)
	o := sha1.New()
	o.Write([]byte(verifyStr))
	signature := hex.EncodeToString(o.Sum(nil))
	response.Success(c, response.GetSuccess, "查询成功", map[string]interface{}{
		"appid":     config.Global.Config.App.AppID,
		"signature": signature,
		"nonceStr":  nonceStr,
		"timestamp": timestamp,
	})
}
