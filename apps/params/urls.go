package params

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine) {
	// 路由组
	group := e.Group("/params")
	{
		group.GET("/querystring", querystringHandler)
		group.GET("/path/:username/:address", pathHandler)
		group.POST("/form", formHandler)
		group.POST("/json", jsonHandler)
		group.GET("/shoudbind", shoudbindHandler)

		// HTTP重定向
		group.GET("/redirect1", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com/")
		})
		// 路由重定向
		group.GET("/redirect2", func(c *gin.Context) {
			c.Request.URL.Path = "/querystring"
			e.HandleContext(c)
		})
	}
}
