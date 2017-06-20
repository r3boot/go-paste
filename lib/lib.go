package lib

import (
	"gopkg.in/redis.v3"
)

var log Log
var Config *AppConfig
var Redis *redis.Client

func Setup(l Log, c *AppConfig) {
	log = l
	Config = c

	Redis = redis.NewClient(&redis.Options{
		Addr:     Config.RedisAddr,
		Password: Config.RedisPasswd,
		DB:       Config.RedisDB,
	})

	SetupServer()
}
