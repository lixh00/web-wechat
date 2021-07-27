package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// Log 日志工具
var Log *zap.SugaredLogger

// InitLogger 初始化日志工具
func InitLogger() {
	// 配置 sugaredLogger
	writer := zapcore.AddSync(os.Stdout)

	// 自定义时间输出格式
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format("2006-01-02 15:04:05.000") + "]")
	}

	// 格式相关的配置
	encoderConfig := zap.NewProductionEncoderConfig()
	// 修改时间戳的格式
	encoderConfig.EncodeTime = customTimeEncoder
	// 日志级别使用大写带颜色显示
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	// 将日志级别设置为 DEBUG
	core := zapcore.NewCore(encoder, writer, zapcore.DebugLevel)
	// 增加 caller 信息
	logger := zap.New(core, zap.AddCaller())
	Log = logger.Sugar()
}
