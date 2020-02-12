package main

import (
	"flag"
	"github.com/r3boot/go-paste/lib"
)

// Default values for commandline paramenters
const (
	defDebugValue     = false
	defDebugHelp      = "Enable debugging output"
	defBindIpValue    = "0.0.0.0"
	defBindIpHelp     = "IP to listen on"
	defBindPortValue  = 8080
	defBindPortHelp  = "Port to listen on"
	defRedisAddressValue = "localhost:6379"
	defRedisAddressHelp  = "Address on which redis is running"
	defRedisPassValue = ""
	defRedisPassHelp  = "Password used to connect to redis"
	defRedisDbValue    = 0
	defRedisDbHelp     = "Redis DB to connect to"
	defTitleValue      = "Pastebin"
	defTitleHelp = "Title to use for the webapp"
)

// Logging framework
var Log lib.Log

func main() {
	var (
		debug     = flag.Bool("debug", defDebugValue, defDebugHelp)
		bindip    = flag.String("listen", defBindIpValue, defBindIpHelp)
		bindport  = flag.Int("port", defBindPortValue, defBindPortHelp)
		redisaddr = flag.String("redis-host", defRedisAddressValue, defRedisAddressHelp)
		redispass = flag.String("redis-pass", defRedisPassValue, defRedisPassHelp)
		redisdb   = flag.Int64("redis-db", defRedisDbValue, defRedisDbHelp)
		title     = flag.String("title", defTitleValue, defTitleHelp)
	)

	// Parse commandline paramenters
	flag.Parse()

	// Configure logging framework
	Log.UseDebug = *debug
	Log.UseVerbose = *debug
	Log.UseTimestamp = true
	Log.Debug("Logging initialized")

	lib.Setup(Log, &lib.AppConfig{
		BindIp:      *bindip,
		BindPort:    *bindport,
		RedisAddr:   *redisaddr,
		RedisPasswd: *redispass,
		RedisDB:     *redisdb,
		Title:       *title,
	})


	lib.RunServer()
}
