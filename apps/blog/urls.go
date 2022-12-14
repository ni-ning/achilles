package blog

import (
	"github.com/gin-gonic/gin"
)

var article = NewArticle()
var tag = NewTag()

func Routers(e *gin.Engine) {
	apiv1 := e.Group("/api/v1")
	{
		apiv1.POST("/tags", tag.Create)       // 新增标签
		apiv1.GET("/tags", tag.List)          // 获取标签列表
		apiv1.PUT("/tags/:id", tag.Update)    // 更新指定标签
		apiv1.DELETE("/tags/:id", tag.Delete) // 删除指定标签

		apiv1.POST("/articles", article.Create)           // 新增文章
		apiv1.DELETE("/articles/:id", article.Delete)     // 删除指定文章
		apiv1.PUT("/articles/:id", article.Update)        // 更新指定文章
		apiv1.PATCH("/articles/:id/state", article.Patch) // 禁用或启用文章
		apiv1.GET("/articles/:id", article.Get)           // 获取指定文章
		apiv1.GET("/articles", article.List)              // 获取文章列表
	}
}
