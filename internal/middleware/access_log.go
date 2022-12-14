package middleware

import (
	"achilles/global"
	"bytes"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AccessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w AccessLogWriter) Write(p []byte) (int, error) {
	// 截流多写一次
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer = bodyWriter重写接入一下
		bodyWriter := &AccessLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyWriter

		beginTime := time.Now().Unix()
		c.Next()
		endTime := time.Now().Unix()

		global.Logger.Info("Access Log",
			zap.String("method", c.Request.Method),
			zap.Int("status_code", bodyWriter.Status()),
			zap.Int64("begin_time", beginTime),
			zap.Int64("endTime", endTime),
			zap.String("request", c.Request.PostForm.Encode()),
			zap.String("response", bodyWriter.body.String()),
		)
	}
}
