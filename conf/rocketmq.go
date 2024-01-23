package conf

type RocketmqConfig struct {
	Host  string `yaml:"host"`
	Port  string `yaml:"port"`
	Retry int    `yaml:"retry"` // 重试次数
	Topic string `yaml:"topic"`
	Group string `yaml:"group"`
}
