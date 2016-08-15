package main

import (
	"flag"
	"github.com/r3boot/go-paste/lib"
	"github.com/r3boot/rlib/logger"
)

// Default values for commandline paramenters
const D_DEBUG bool = false
const D_BINDIP string = "0.0.0.0"
const D_BINDPORT int = 8080
const D_REDISADDR string = "localhost:6379"
const D_REDISPASS string = ""
const D_REDISDB int64 = 0
const D_TEMPLATE string = "/usr/share/go-paste.html"
const D_TITLE string = "Pastebin"

// Commandline paramenters
var debug = flag.Bool("D", D_DEBUG, "Enable debug output")
var bindip = flag.String("l", D_BINDIP, "IP to listen on")
var bindport = flag.Int("p", D_BINDPORT, "Port to listen on")
var template = flag.String("t", D_TEMPLATE, "Template to use")
var redisaddr = flag.String("ra", D_REDISADDR, "Address on which redis is running")
var redispass = flag.String("rp", D_REDISPASS, "Password used to connect to redis")
var redisdb = flag.Int64("rd", D_REDISDB, "Redis DB to connect to")
var title = flag.String("T", D_TITLE, "Title to use for the webapp")

// Logging framework
var Log logger.Log

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
		Template:    *template,
	})
}

func main() {
	lib.RunServer()
}
