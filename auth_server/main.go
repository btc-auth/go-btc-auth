package main

import (
	"auth_server/libs"
	"auth_server/views"

	"github.com/gin-gonic/gin"
)

func main() {
	config := map[string]string{
		"port":              "8085",
		"host_url":          "http://192.168.1.3",
		"token_timeout_sec": "60",
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	tokenManager := libs.InitTokenManager()

	// 将 tokenManager 传递给每个页面
	r.Use(func(c *gin.Context) {
		c.Set("tokenManager", tokenManager)
		c.Set("config", config)
		c.Next()
	})

	// 将authMiddleware添加到全局中间件
	r.Use(libs.AuthMiddleware())

	r.GET("/", views.HomePage)
	r.GET("/auth", views.AuthPage)
	r.GET("/result", views.ResultPage)
	r.GET("/dashboard", views.DashboardPage)
	r.GET("/logout", views.LogoutPage)

	r.Run(":" + config["port"])
}
