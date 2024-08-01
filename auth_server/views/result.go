package views

import (
	"auth_server/libs"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

// 返回的数据结构
type ResultPageData struct {
	Status       string `json:"status"`                    // 认证状态 success/timeout/unmatched
	Message      string `json:"message"`                   // 错误信息，或者成功信息
	Redirect     string `json:"redirect,omitempty"`        // success 下使用。不是必须的，成功的时候，给出跳转 url
	NewToken     string `json:"new_token,omitempty"`       // failure 下使用。不是必须的，失败的时候，给出新的 token
	NewQrUrl     string `json:"new_qr_url,omitempty"`      // failure 下使用。不是
	NewQrDataUri string `json:"new_qr_data_uri,omitempty"` // failure 下使用。不是
}

/*
ResultPageData 有三种情况
success: 认证成功
	status: success
	message: authentication successful
	redirect: 跳转的 url

timeout: 超时
	status: timeout
	message: token timeout
	new_token: 新的 token
	new_qr_url: 新的二维码 url
	new_qr_data_uri: 新的二维码 Base64 编码

unmatched: 未匹配
	status: failure
	message: no authentication
*/

func ResultPage(c *gin.Context) {
	config := c.MustGet("config").(map[string]string)
	tokenManager := c.MustGet("tokenManager").(*libs.TokenManager)
	token := c.Query("token")

	token_obj, err := tokenManager.CheckToken(token)
	if err != nil {
		// 超时，更新 token
		token := tokenManager.NewToken()
		qrCodeURL := libs.GetAuthUrl(config["host_url"], config["port"], token)
		// 生成二维码
		var png []byte
		png, thiserr := qrcode.Encode(qrCodeURL, qrcode.Medium, 256)
		if thiserr != nil {
			fmt.Println("生成二维码失败")
			c.String(http.StatusInternalServerError, "生成二维码失败")
			return
		}

		// 将二维码转换为Base64编码
		base64QRCode := base64.StdEncoding.EncodeToString(png)
		qrCodeDataURI := "data:image/png;base64," + base64QRCode

		c.JSON(http.StatusOK, ResultPageData{
			Status:       "timeout",
			Message:      "token timeout",
			NewToken:     token,
			NewQrUrl:     qrCodeURL,
			NewQrDataUri: qrCodeDataURI,
		})
		return
	} else {
		if token_obj.PublicKey == "" {
			// 没有匹配
			c.JSON(http.StatusOK, ResultPageData{
				Status:  "unmatched",
				Message: "no authentication",
			})
			return
		} else {
			// 匹配上了
			c.SetCookie("session_token", token, 3600, "/", "", false, true)
			c.JSON(http.StatusOK, ResultPageData{
				Status:   "success",
				Message:  "authentication successful",
				Redirect: libs.GetLoginHomeUrl(config["host_url"], config["port"]),
			})
			return
		}
	}
}
