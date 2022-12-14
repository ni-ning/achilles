package orm

import "github.com/gin-gonic/gin"

func Routers(e *gin.Engine) {
	group := e.Group("/orm")

	group.GET("/create", createHandler)
	group.GET("/select", selectHandler)
}
