package logger

import (
	"achilles/pkg/setting"
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(appSetting *setting.AppSettingS) *zap.Logger {
	// 1. Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 时间格式 2022-12-08T18:24:07.979+0800
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 大写 "level":"INFO"
	encoderConfig.TimeKey = "timestamp"                     // "timestamp":"2022-12-08T18:41:35.596+0800"
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 2. WriterSyncer 同时输出到 lumberJackLogger 和控制台
	logLocation := fmt.Sprintf("%s%s%s",
		appSetting.LogSavePath,
		appSetting.LogFileName,
		appSetting.LogFileExt)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logLocation, // 日志文件的位置
		MaxSize:    1,           // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 5,           // 保留旧文件的最大个数
		MaxAge:     30,          // 保留旧文件的最大天数
		Compress:   false,       // 是否压缩/归档旧文件
	}
	writeSyncer := zapcore.AddSync(io.MultiWriter(lumberJackLogger, os.Stdout))

	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	// 调用函数信息 如 "caller":"day15/main.go:27", zap.AddCallerSkip(1)用于额外封装一层场景
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}
