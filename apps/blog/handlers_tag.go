package blog

import (
	"achilles/pkg/app"
	"achilles/pkg/convert"
	"achilles/pkg/errcode"

	"github.com/gin-gonic/gin"
)

type Tag struct{}

func NewTag() Tag {
	return Tag{}
}

func (t Tag) Get(c *gin.Context) {}

// @Summary 获取多个标签
// @Tags 博客后台
// @Produce  json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param name query string false "标签名"
// @Param state query int false "标签状态"
// @Success 200 {object} model.Tag "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags [get]
func (t Tag) List(c *gin.Context) {
	req := TagListCountRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	pageSize, page := app.GetPageSize(c), app.GetPage(c)

	tags, err := GetTagList(req, page, pageSize)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorGetTagListFail)
		return
	}
	totalRows, err := GetTagCount(req)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorCountTagFail)
		return
	}

	response.ToResponseList(tags, totalRows)
}

// @Summary 新增标签
// @Tags 博客后台
// @Produce  json
// @Param tag body TagCreateRequest true "标签名"
// @Success 200 {object} model.Tag "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags [post]
func (t Tag) Create(c *gin.Context) {
	req := TagCreateRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := CreateTag(req)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorCreateTagFail)
		return
	}
	response.ToResponse(nil)
}

// @Summary 更新指定标签
// @Tags 博客后台
// @Produce  json
// @Param id path int true "标签 ID"
// @Param tag body TagUpdateRequest true "标签信息"
// @Success 200 {object} model.Tag "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags/{id} [put]
func (t Tag) Update(c *gin.Context) {
	// 厉害的了，初始化带过来，但是得注意 body中的id优先级更高
	req := TagUpdateRequest{ID: convert.StrTo(c.Param("id")).MustUInt64()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := UpdateTag(req)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorUpdateTagFail)
		return
	}
	response.ToResponse(nil)
}

// @Summary 删除指定标签
// @Tags 博客后台
// @Produce  json
// @Param id path int true "标签 ID"
// @Success 200 {object} model.Tag "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags/{id} [delete]
func (t Tag) Delete(c *gin.Context) {
	req := TagDeleteRequest{ID: convert.StrTo(c.Param("id")).MustUInt64()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := DeleteTag(req)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorDeleteTagFail)
		return
	}
	response.ToResponse(nil)
}
