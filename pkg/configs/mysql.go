package configs

type MysqlConfig struct {
	Host             string `yaml:"host"`
	Port             string `yaml:"port"`
	Database         string `yaml:"database"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	Charset          string `yaml:"charset"`
	Timeout          string `yaml:"timeout"`
	TablePrefix      string `yaml:"table_prefix"`
	LogMode          int    `yaml:"log_mode"`
	MaxOpenConns     int    `yaml:"max_open_conns"`
	MaxIdleConns     int    `yaml:"max_idle_conns"`
	MaxLifetimeConns int64  `yaml:"max_lifetime_conns"`
}
