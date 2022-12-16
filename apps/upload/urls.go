package upload

import (
	"achilles/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routers(e *gin.Engine) {
	u := NewUpload()
	e.POST("/upload/file", u.UploadFile) // 上传文件
	e.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))
}
