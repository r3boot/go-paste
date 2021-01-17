package lib

import (
	"fmt"
	"os"
	"time"
)

const (
	infoMessage    byte = 0x0
	warningMessage byte = 0x1
	fatalMessage   byte = 0x2
	verboseMessage byte = 0x3
	debugMessage   byte = 0x4
)

var messageString = map[byte]string{
	infoMessage:    "I",
	warningMessage: "W",
	fatalMessage:   "F",
	verboseMessage: "V",
	debugMessage:   "D",
}

type Log struct {
	UseDebug        bool
	UseVerbose      bool
	UseTimestamp    bool
	TimestampFormat string
	TestFd          *os.File
}

func (l Log) Message(log_level byte, message ...interface{}) {
	var (
		msg string
		fd  *os.File
	)

	if l.UseTimestamp {
		if len(l.TimestampFormat) == 0 {
			l.TimestampFormat = time.RFC3339
		}
		timestamp := time.Now().Format(time.RFC3339)
		msg = timestamp + " " + messageString[log_level] + ":"
	} else {
		msg = messageString[log_level] + ":"
	}

	all := append([]interface{}{msg}, message...)

	if log_level == fatalMessage {
		fd = os.Stderr
	} else {
		fd = os.Stdout
	}

	if l.TestFd != nil {
		fd = l.TestFd
	}

	fmt.Fprintln(fd, all...)
}

func (l Log) Info(message ...interface{}) {
	l.Message(infoMessage, message...)
}

func (l Log) Warn(message ...interface{}) {
	l.Message(warningMessage, message...)
}

func (l Log) Error(message ...interface{}) {
	l.Message(fatalMessage, message...)
	if l.TestFd == nil {
		os.Exit(1)
	}
}

func (l Log) Verbose(message ...interface{}) {
	if l.UseDebug || l.UseVerbose {
		l.Message(verboseMessage, message...)
	}
}

func (l Log) Debug(message ...interface{}) {
	if l.UseDebug {
		l.Message(debugMessage, message...)
	}
}
