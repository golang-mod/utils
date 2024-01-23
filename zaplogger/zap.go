package zaplogger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Filename   string `json:"filename" yaml:"filename"`       //日志文件存放目录
	MaxSize    int    `json:"max_size" yaml:"max_size"`       //文件大小限制,单位MB
	MaxAge     int    `json:"max_age" yaml:"max_age"`         //日志文件保留天数
	MaxBackups int    `json:"max_backups" yaml:"max_backups"` //最大保留日志文件数量
	LocalTime  bool   `json:"localtime" yaml:"localtime"`     // 本地时区
	Compress   bool   `json:"compress" yaml:"compress"`       //是否压缩处理
	Level      string `json:"level" yaml:"level"`             //日志级别
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(time.DateTime))
}

// New 日志
// logger.Error("无法获取网址",
//
//	zap.String("data", "http://www.baidu.com"),
//	zap.Int("attempt", 3),
//	zap.Duration("backoff", time.Second),
//
// )
// logger.Sugar().Debug("test1280's debug")
// logger.Sugar().Infof("test1280's %s", "infof")
// logger.Sugar().Warnf("test1280's %s", "warnf")
// logger.Sugar().Error("test1280's error")
func New(config Config, datatype string) *zap.Logger {
	atomicLevel := zap.NewAtomicLevel()
	switch config.Level {
	case "DEBUG":
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case "INFO":
		atomicLevel.SetLevel(zapcore.InfoLevel)
	case "WARN":
		atomicLevel.SetLevel(zapcore.WarnLevel)
	case "ERROR":
		atomicLevel.SetLevel(zapcore.ErrorLevel)
	case "DPANIC":
		atomicLevel.SetLevel(zapcore.DPanicLevel)
	case "PANIC":
		atomicLevel.SetLevel(zapcore.PanicLevel)
	case "FATAL":
		atomicLevel.SetLevel(zapcore.FatalLevel)
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "name",
		CallerKey:     "line",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
	}
	var enc zapcore.Encoder
	//文件writeSyncer
	writer := WriteSyncer(config)
	if datatype == "JSON" {
		encoderConfig.LineEnding = zapcore.DefaultLineEnding
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		encoderConfig.EncodeTime = customTimeEncoder
		encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
		encoderConfig.EncodeName = zapcore.FullNameEncoder
		enc = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		enc = zapcore.NewConsoleEncoder(encoderConfig)
	}
	zapCore := zapcore.NewCore(
		enc,
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(writer),
			zapcore.AddSync(os.Stdout),
		),
		atomicLevel,
	)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 设置初始化字段
	//filed := zap.Fields(zap.String("serviceName", "serviceName"))
	logger := zap.New(zapCore, caller)

	defer logger.Sync()
	return logger
}

func WriteSyncer(config Config) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   config.Filename,   //日志文件存放目录
		MaxSize:    config.MaxSize,    //文件大小限制,单位MB
		MaxBackups: config.MaxBackups, //最大保留日志文件数量
		MaxAge:     config.MaxAge,     //日志文件保留天数
		LocalTime:  config.LocalTime,  //本地时区
		Compress:   config.Compress,   //是否压缩处理
	}
}
