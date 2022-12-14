package params

import (
	"achilles/global"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func querystringHandler(c *gin.Context) {
	username := c.DefaultQuery("username", "linda") // 默认 "linda"
	address := c.Query("address")                   // 默认零值

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"address":  address,
	})
}

func formHandler(c *gin.Context) {
	// DefaultPostForm取不到值时会返回指定的默认值
	username := c.DefaultPostForm("username", "linda") // 默认 "linda"
	address := c.PostForm("address")                   // 默认零值

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"address":  address,
	})
}

func jsonHandler(c *gin.Context) {
	// 注意：下面为了举例子方便，暂时忽略了错误处理
	b, _ := c.GetRawData() // 从c.Request.Body读取请求数据

	// 定义map或结构体
	var m map[string]interface{}
	// 反序列化
	_ = json.Unmarshal(b, &m)

	fmt.Printf("type: %T\nvalue:%#v\n", m, m)
	global.Logger.Info("JSON数据", zap.Any("m", m))

	c.JSON(http.StatusOK, m)
}

func pathHandler(c *gin.Context) {
	username := c.Param("username")
	address := c.Param("address")

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"address":  address,
	})
}

type Login struct {
	// form:"user" 校验 FormPost提交数据
	// json:"user" 校验json提交数据
	// binding:"required"  如果没有改选项，默认零值
	Uesr     string `form:"user" json:"user" binding:"required"`
	Password string `json:"password"`
}

func shoudbindHandler(c *gin.Context) {

	var login Login
	// ShouldBind()会根据请求的Content-Type自行选择绑定器
	if err := c.ShouldBind(&login); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "参数错误",
		})

		global.Logger.Error("参数错误",
			zap.String("path", c.Request.URL.Path),
			zap.String("error", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": login.Uesr,
		"password": login.Password,
	})
}
