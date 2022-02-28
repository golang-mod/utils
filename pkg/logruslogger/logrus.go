package logruslogger

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func New(logType string) *logrus.Logger {
	//文件writeSyncer
	writer := &lumberjack.Logger{
		Filename:   "./storage/logs/" + logType + ".log", //日志文件存放目录
		MaxSize:    1,                                    //文件大小限制,单位MB
		MaxBackups: 5,                                    //最大保留日志文件数量
		MaxAge:     30,                                   //日志文件保留天数
		LocalTime:  true,                                 // 本地时区
		Compress:   false,                                //是否压缩处理
	}

	//实例化
	log := logrus.New()

	//设置输出
	//log.Out = writer
	log.SetOutput(writer)
	log.SetOutput(os.Stdout)
	//设置日志级别
	log.SetLevel(logrus.DebugLevel)

	//设置日志格式
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return log
}
