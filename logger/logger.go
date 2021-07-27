package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Log 日志工具
var Log *zap.SugaredLogger

// InitLogger 初始化日志工具
func InitLogger() {
	// 配置 sugaredLogger
	writer := zapcore.AddSync(os.Stdout)

	// 格式相关的配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 修改时间戳的格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 日志级别使用大写显示
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, writer, zapcore.DebugLevel) // 将日志级别设置为 DEBUG
	logger := zap.New(core, zap.AddCaller())                     // 增加 caller 信息
	Log = logger.Sugar()
}
