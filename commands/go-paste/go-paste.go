package main

import (
	"flag"
	"github.com/r3boot/go-paste/lib"
)

// Default values for commandline paramenters
const (
	D_DEBUG     bool   = false
	D_BINDIP    string = "0.0.0.0"
	D_BINDPORT  int    = 8080
	D_REDISADDR string = "localhost:6379"
	D_REDISPASS string = ""
	D_REDISDB   int64  = 0
	D_TITLE     string = "Pastebin"
)

// Commandline paramenters
var (
	debug     = flag.Bool("D", D_DEBUG, "Enable debug output")
	bindip    = flag.String("l", D_BINDIP, "IP to listen on")
	bindport  = flag.Int("p", D_BINDPORT, "Port to listen on")
	redisaddr = flag.String("ra", D_REDISADDR, "Address on which redis is running")
	redispass = flag.String("rp", D_REDISPASS, "Password used to connect to redis")
	redisdb   = flag.Int64("rd", D_REDISDB, "Redis DB to connect to")
	title     = flag.String("T", D_TITLE, "Title to use for the webapp")
)

// Logging framework
var Log lib.Log

func init() {
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
}

func main() {
	lib.RunServer()
}
