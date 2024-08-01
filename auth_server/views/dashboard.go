package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DashboardPage(c *gin.Context) {
	publickeyHex := c.GetString("publickey_hex")
	if publickeyHex == "" {
		c.Redirect(http.StatusFound, "/")
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"publickey_hex": publickeyHex,
	})
}
