package todo

import (
	"achilles/global"
	"achilles/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// getTodoHandler 获取全部清单列表
func getTodoHandler(c *gin.Context) {
	var todos []model.Todo
	if err := global.DBEngine.Debug().Find(&todos).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  err.Error,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": todos,
	})
}

// postTodoHandler 创建清单
func postTodoHandler(c *gin.Context) {
	var todo model.Todo
	if err := c.ShouldBind(&todo); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "参数错误1",
		})
		return
	}
	if todo.Title == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "参数错误2",
		})
		return
	}
	if err := global.DBEngine.Create(&todo).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "创建错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
}

// putTodoHandler 修改清单
func putTodoHandler(c *gin.Context) {
	var todo model.Todo
	if err := c.ShouldBind(&todo); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "参数错误1",
		})
		return
	}
	if err := global.DBEngine.First(&todo, todo.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "数据不存在",
		})
		return
	}
	if err := global.DBEngine.Model(&todo).Update("status", todo.Status).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
}

// deleteTodoHandler 删除清单
func deleteTodoHandler(c *gin.Context) {
	var todo model.Todo
	if err := c.ShouldBind(&todo); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "参数错误1",
		})
		return
	}
	if err := global.DBEngine.First(&todo, todo.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "数据不存在",
		})
		return
	}
	if err := global.DBEngine.Model(&todo).Update("is_deleted", true).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
}
