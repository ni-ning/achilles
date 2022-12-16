package upload

import (
	"achilles/global"
	"achilles/pkg/app"
	"achilles/pkg/convert"
	"achilles/pkg/errcode"
	"achilles/pkg/upload"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Upload struct{}

func NewUpload() Upload {
	return Upload{}
}

// @Summary 文件上传
// @Tags 文件上传
// @Accept	multipart/form-data
// @Accept	application/x-www-form-urlencoded
// @Produce json
// @Param file formData file true "文件"
// @Param type query int true "文件类型"
// @Success 200 "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /upload/file [post]
func (u Upload) UploadFile(c *gin.Context) {
	response := app.NewResponse(c)
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	fileType := convert.StrTo(c.DefaultQuery("type", "1")).MustInt()
	if fileHeader == nil || fileType <= 0 {
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	fileInfo, err := UploadFile(upload.FileType(fileType), file, fileHeader)
	if err != nil {
		global.Logger.Error("UploadFile Error", zap.String("error", err.Error()))
		response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
		return
	}

	response.ToResponse(gin.H{
		"file_access_url": fileInfo.AccessUrl,
	})
}
