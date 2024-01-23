package redisx

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

type Config struct {
	Host            string `yaml:"host"`
	Password        string `yaml:"password"`
	Port            string `yaml:"port"`
	DefaultDatabase int    `yaml:"default_database"`
}

type Redisx struct {
	config Config
}

func New(config Config) *Redisx {
	x := Redisx{}
	x.config = config
	return &x
}
func (x *Redisx) Client() (*redis.Client, error) {
	return x.connection()
}
func (x *Redisx) connection() (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{ // 连接服务
		Addr:     x.config.Host + ":" + x.config.Port, // string
		Password: x.config.Password,                   // string
		DB:       x.config.DefaultDatabase,            // int
	})
	var ping string
	ping, err = client.Ping(context.Background()).Result() // 心跳
	if err != nil {
		log.Println("redis", ping)
		return
	}
	return
}
