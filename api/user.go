package api

import (
	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
	"main.go/service"
)

func UserRegister(c *gin.Context) {
	var userRegisterService service.UserRegisterService
	if err := c.ShouldBind(&userRegisterService); err == nil {
		userRegisterService.UserName = c.PostForm("username")
		userRegisterService.PassWord = c.PostForm("password")
		res := userRegisterService.Register()
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		logging.Info(err)
	}

}
