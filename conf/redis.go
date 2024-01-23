package conf

type RedisConfig struct {
	Host            string `yaml:"host"`
	Password        string `yaml:"password"`
	Port            string `yaml:"port"`
	DefaultDatabase int    `yaml:"default_database"`
}
