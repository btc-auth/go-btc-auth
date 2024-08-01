package views

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"

	"auth_server/libs"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

func HomePage(c *gin.Context) {
	publickeyHex := c.GetString("publickey_hex")
	if publickeyHex != "" {
		c.Redirect(http.StatusFound, "/dashboard")
	}
	tokenManager := c.MustGet("tokenManager").(*libs.TokenManager)
	config := c.MustGet("config").(map[string]string)
	token := tokenManager.NewToken()

	qrCodeURL := libs.GetAuthUrl(config["host_url"], config["port"], token)

	// 生成二维码
	var png []byte
	png, err := qrcode.Encode(qrCodeURL, qrcode.Medium, 256)
	if err != nil {
		fmt.Println("生成二维码失败")
		c.String(http.StatusInternalServerError, "生成二维码失败")
		return
	}

	// 将二维码转换为Base64编码
	base64QRCode := base64.StdEncoding.EncodeToString(png)
	qrCodeDataURI := "data:image/png;base64," + base64QRCode

	fmt.Println("二维码地址:", qrCodeURL)
	fmt.Println("二维码Base64编码:", qrCodeDataURI)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"qrCodeDataURI": template.URL(qrCodeDataURI),
		"qrCodeURL":     qrCodeURL,
		"token":         token,
	})
}
