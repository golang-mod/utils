package configs

type LoggerConfig struct {
	Filename   string `json:"filename" yaml:"filename"`       //日志文件存放目录
	MaxSize    int    `json:"max_size" yaml:"max_size"`       //文件大小限制,单位MB
	MaxAge     int    `json:"max_age" yaml:"max_age"`         //日志文件保留天数
	MaxBackups int    `json:"max_backups" yaml:"max_backups"` //最大保留日志文件数量
	LocalTime  bool   `json:"localtime" yaml:"localtime"`     // 本地时区
	Compress   bool   `json:"compress" yaml:"compress"`       //是否压缩处理
	Level      string `json:"level" yaml:"level"`             //日志级别
}
