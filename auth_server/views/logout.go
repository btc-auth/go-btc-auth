package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogoutPage(c *gin.Context) {
	// 清除 session_token
	c.SetCookie("session_token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}
