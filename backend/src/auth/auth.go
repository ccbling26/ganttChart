package auth

import (
	"backend/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type DataGetTicket struct {
	Ticket   string `json:"ticket"`
	ExpireIn int    `json:"expire_in"`
}

type ResponseGetTicket struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data DataGetTicket `json:"data"`
}

type ResponseAuthWithTenantAccessToken struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

func GetTicket() (string, error) {
	// https://open.feishu.cn/document/ukTMukTMukTM/uYTM5UjL2ETO14iNxkTN/h5_js_sdk/authorization
	appConfig := config.Global.Config.App
	tenantAccessToken, err := authWithTenantAccessToken()
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s%s", appConfig.Host, config.JSAPI_TICKET_URI)
	reader := bytes.NewReader([]byte{})
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return "", err
	}
	request.Header.Set("Authorization", "Bearer "+tenantAccessToken)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if response.StatusCode != 200 {
		return "", errors.New("请求 ticket 失败，状态码为：" + strconv.Itoa(response.StatusCode))
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var res ResponseGetTicket
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}
	if res.Code != 0 {
		return "", errors.New(res.Msg)
	}
	return res.Data.Ticket, nil
}

func authWithTenantAccessToken() (string, error) {
	// https://open.feishu.cn/document/server-docs/authentication-management/access-token/tenant_access_token_internal
	appConfig := config.Global.Config.App
	url := fmt.Sprintf("%s%s", appConfig.Host, config.TENANT_ACCESS_TOKEN_URI)
	reqBody := map[string]string{
		"app_id":     appConfig.AppID,
		"app_secret": appConfig.AppSecret,
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	reader := bytes.NewReader(data)
	response, err := http.Post(url, "application/json; charset=utf-8", reader)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if response.StatusCode != 200 {
		return "", errors.New("请求 tenant_access_token 失败，状态码为：" + strconv.Itoa(response.StatusCode))
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var res ResponseAuthWithTenantAccessToken
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}
	if res.Code != 0 {
		return "", errors.New(res.Msg)
	}
	return res.TenantAccessToken, nil
}
