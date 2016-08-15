package main

import (
	"flag"
	"github.com/r3boot/go-paste/lib"
	"github.com/r3boot/rlib/logger"
	_ "time"
)

// Default values for commandline paramenters
const D_DEBUG bool = false

// Commandline paramenters
var debug = flag.Bool("D", D_DEBUG, "Enable debug output")

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
		BindIp:      "0.0.0.0",
		BindPort:    8080,
		RedisAddr:   "redis.service.local:6379",
		RedisPasswd: "",
		RedisDB:     0,
		Title:       "as65342 paste",
		Template:    "./templates/index.html",
	})
}

func main() {
	lib.RunServer()
}
