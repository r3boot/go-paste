package lib

import (
	"github.com/r3boot/rlib/logger"
	"gopkg.in/redis.v3"
)

var Log logger.Log
var Config *AppConfig
var Redis *redis.Client

func Setup(l logger.Log, c *AppConfig) {
	Log = l
	Config = c

	Redis = redis.NewClient(&redis.Options{
		Addr:     Config.RedisAddr,
		Password: Config.RedisPasswd,
		DB:       Config.RedisDB,
	})

	SetupServer()
}
