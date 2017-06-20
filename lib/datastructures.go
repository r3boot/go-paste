package lib

import (
	"time"
)

type AppConfig struct {
	BindIp      string
	BindPort    int
	RedisAddr   string
	RedisPasswd string
	RedisDB     int64
	Title       string
}

type Paste struct {
	Hash       string
	Expiration time.Duration
	Content    []byte
}

type NewPaste struct {
	Title         string
	ExpiryOptions map[string]string
}
