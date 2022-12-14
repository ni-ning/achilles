package todo

import "github.com/gin-gonic/gin"

func Routers(e *gin.Engine) {
	g := e.Group("/todo")
	{
		g.GET("/todo", getTodoHandler)
		g.POST("/todo", postTodoHandler)
		g.PUT("/todo", putTodoHandler)
		g.DELETE("/todo", deleteTodoHandler)
	}
}
