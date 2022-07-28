package service

import (
	"main.go/model"
	"main.go/serializer"
)

type UserRegisterService struct {
	UserName string `json:"username" from:"username"`
	PassWord string `json:"password" from:"password"`
}

func (service *UserRegisterService) Register() serializer.Response {
	var user model.User
	count := 0
	model.DB.Model(&model.User{}).Where("username=?", service.UserName).First(&user).Count(&count)
	if count != 0 {
		// 已经有人注册
		return serializer.Response{
			Status: 400,
			Msg:    "用户名已注册",
		}
	}
	user = model.User{
		UserName: service.UserName,
	}
	if err := user.SetPassword(service.PassWord); err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "密码加密出错",
		}
	}
	model.DB.Create(&user)
	return serializer.Response{
		Status: 200,
		Msg:    "用户创建成功",
	}
}
