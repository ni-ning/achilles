package middleware

import (
	"achilles/apps/auth"
	"achilles/pkg/app"
	"achilles/pkg/errcode"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		//@Param token header string true "Token"
		var (
			token string
			ecode = errcode.Success
		)
		if s, exist := c.GetQuery("token"); exist {
			token = s
		} else {
			token = c.GetHeader("token")
		}

		response := app.NewResponse(c)
		if token == "" {
			response.ToErrorResponse(errcode.UnauthorizedTokenLack)
			c.Abort()
			return
		}

		claims, err := app.ParseToken(token)
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				ecode = errcode.UnauthorizedTokenTimeout
			default:
				ecode = errcode.UnauthorizedTokenError
			}
			response.ToErrorResponse(ecode)
			c.Abort()
			return
		}

		existed, account := auth.GetAccountById(claims.AccountID)
		if existed {
			c.Set("account", account)
		}

		c.Next()
	}
}
