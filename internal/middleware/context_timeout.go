package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func ContextTimeout(t time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		// 接过来
		ctx, cancel := context.WithTimeout(c.Request.Context(), t)
		defer cancel()
		// 传过去
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
