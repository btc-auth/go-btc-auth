package views

import (
	"net/http"

	"auth_server/libs"

	"github.com/gin-gonic/gin"
)

// 返回的数据结构
type AuthPageData struct {
	Status   string `json:"status"`             // 认证状态 success/timeout/failure
	Message  string `json:"message"`            // 错误信息，或者成功信息
	Redirect string `json:"redirect,omitempty"` // success 下使用。不是必须的，成功的时候，给出跳转 url
}

func AuthPage(c *gin.Context) {
	config := c.MustGet("config").(map[string]string)
	tokenManager := c.MustGet("tokenManager").(*libs.TokenManager)

	token := c.Query("token")
	sign := c.Query("sign")
	publicKeyHex := c.Query("public_key_hex")

	token_obj, err := tokenManager.CheckToken(token)

	if err != nil {
		c.JSON(http.StatusBadRequest, AuthPageData{
			Status:  "timeout",
			Message: "token timeout",
		})
		return
	}

	qrCodeURL := libs.GetAuthUrl(config["host_url"], config["port"], token)
	valid, err := libs.VerifySignature(qrCodeURL, sign, publicKeyHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, AuthPageData{
			Status:  "failure",
			Message: err.Error(),
		})
		return
	}

	if valid {
		token_obj.Auth(publicKeyHex)
		// 签名验证成功
		c.JSON(http.StatusOK, AuthPageData{
			Status:   "success",
			Message:  "Authentication successful",
			Redirect: libs.GetLoginHomeUrl(config["host_url"], config["port"]),
		})
	} else {
		// 签名验证失败
		c.JSON(http.StatusUnauthorized, AuthPageData{
			Status:  "failure",
			Message: "Invalid signature",
		})
	}

}
