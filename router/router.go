package router

import (
	"github.com/gin-gonic/gin"
	"main.go/api"
	"main.go/service"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("html/*")
	r.Use(gin.Recovery(), gin.Logger()) //用于恢复没有考虑到的恐慌， 并写入日志
	v1 := r.Group("/")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "success")
		})
		v1.GET("user/register", func(ctx *gin.Context) {
			ctx.HTML(200, "register.html", nil)
		})
		v1.POST("user/register", api.UserRegister)
		v1.GET("ws", service.WsHandler)

	}
	return r
}
