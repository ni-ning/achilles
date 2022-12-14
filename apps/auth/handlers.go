package auth

import (
	"achilles/pkg/app"
	"achilles/pkg/errcode"

	"github.com/gin-gonic/gin"
)

type Account struct{}

func NewAccount() Account {
	return Account{}
}

// @Summary 用户注册
// @Tags 用户管理
// @Produce  json
// @Param tag body AuthRegisterRequest true "用户信息"
// @Success 200 {object} model.Account "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/register [post]
func (a Account) Register(c *gin.Context) {
	req := AuthRegisterRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := AccountRegister(&req)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorRegisterFail)
		return
	}
	response.ToResponse(nil)
}

// @Summary 用户登录
// @Tags 用户管理
// @Produce  json
// @Param tag body AuthLoginRequest true "登录信息"
// @Success 200 {object} model.Account "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/login [post]
func (a Account) Login(c *gin.Context) {
	req := AuthLoginRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	existed, account := AccountLogin(&req)
	if !existed {
		response.ToErrorResponse(errcode.ErrorLoginFail)
		return
	}

	token, err := app.GenerateToken(int64(account.ID))
	if err != nil {
		response.ToErrorResponse(errcode.ErrorGenTokenFail)
		return
	}
	response.ToResponse(map[string]interface{}{"token": token})
}
