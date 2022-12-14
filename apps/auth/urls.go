package auth

import (
	"github.com/gin-gonic/gin"
)

var account = NewAccount()

func Routers(e *gin.Engine) {
	apiv1 := e.Group("/api/v1")
	{
		apiv1.POST("/register", account.Register) // 用户注册
		apiv1.POST("/login", account.Login)       // 用户登录
	}
}
