## Introduction
Go-paste is a small Go-based application which functions as a pastebin. It is designed to be as simple as possible, and consists of a single binary together with a HTML template, and redis to store the paste data in. It features configurable paste expiry based on [time/duration](https://golang.org/pkg/time/).

## Installation
```
export GOPATH=/some/absolute/path
go get github.com/r3boot/go-paste
cd ${GOPATH}/src/github.com/r3boot/go-paste
make
make install
```

## Usage
```
Usage of go-paste:
  -D	Enable debug output
  -T string
    	Title to use for the webapp (default "Pastebin")
  -l string
    	IP to listen on (default "0.0.0.0")
  -p int
    	Port to listen on (default 8080)
  -ra string
    	Address on which redis is running (default "localhost:6379")
  -rd int
    	Redis DB to connect to
  -rp string
    	Password used to connect to redis
  -t string
    	Template to use (default "/usr/share/go-paste.html")
```

## Demo
A live demo of this software, running in a rkt container, can be seen on the URL below:

[paste.as65342.net](https://paste.as65342.net)

## Credits
This pastebin is heavily inspired by [maze.io/paste](https://maze.io/paste), hosted by [tehmaze](https://github.com/tehmaze).
