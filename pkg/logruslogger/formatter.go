package logruslogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

type ExFormatter struct {
}

// Format 自定义格式
// [2022-01-09 15:49:17] [info] 404 | 213.566µs | 127.0.0.1 | GET | /images/icons/gear.png | info log
func (m *ExFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string

	dataType, _ := json.Marshal(entry.Data)
	dataString := string(dataType)

	//HasCaller()为true才会有调用信息
	//entry.Caller.File：文件名
	//entry.Caller.Line: 行号
	//entry.Caller.Function：函数名
	//entry.Caller中还有调用栈相关信息，有需要可以在日志中加入
	//entry.HasCaller() 的判断是必须的，否则如果外部没有设置logrus.SetReportCaller(true)，entry.Caller.*的调用会引发Panic
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		newLog = fmt.Sprintf("[%s] [%s] %s:%d %s | %s | %s\n",
			timestamp, entry.Level, fName, entry.Caller.Line, entry.Caller.Function, entry.Message, dataString)
	} else {
		newLog = fmt.Sprintf("[%s] [%s] %s | %s\n", timestamp, entry.Level, entry.Message, dataString)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}

type AccessFormatter struct {
}

// Format 自定义格式
// [2022-01-09 15:49:17] [info] 404 | 213.566µs | 127.0.0.1 | GET | /images/icons/gear.png | info log
func (m *AccessFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string
	// Gin 参考默认格式
	//[GIN] 2022/01/08 - 23:45:54 | 404 |      78.202µs |       127.0.0.1 | GET      "/images/icons/gear.png"
	newLog = fmt.Sprintf("[%s] [%s] %d | %s | %s | %s | %s \n",
		timestamp,
		entry.Level,
		entry.Data["status_code"],
		entry.Data["latency_time"],
		entry.Data["client_ip"],
		entry.Data["req_method"],
		entry.Data["req_uri"],
	)
	b.WriteString(newLog)
	return b.Bytes(), nil
}
