package upload

import (
	"achilles/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine) {
	u := NewUpload()
	apiv1 := e.Group("/api/v1")
	{
		apiv1.POST("/upload/file", u.UploadFile) // 上传文件
		apiv1.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))
	}
}
