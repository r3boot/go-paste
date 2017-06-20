## Introduction
Go-paste is a small Go-based application which functions as a pastebin. It is designed to be as simple as possible, and consists of a single binary together with a HTML template, and redis to store the paste data in. It features configurable paste expiry based on [time/duration](https://golang.org/pkg/time/).

## Building
```
export GOPATH=/some/absolute/path
go get github.com/r3boot/go-paste
cd ${GOPATH}/src/github.com/r3boot/go-paste
make
```

## Installation
To install go-paste and the gp utility locally, enter the following commands:
```
make install
```

## Usage
The server side can be configured using commandline parameters, as described below:
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
```

The client either reads from stdin or the first file specified on the commandline, and posts the content of that file to a pastebin of your choice:

```
$ gp -h
Usage of gp:
  -D	Enable debug output
  -e string
    	When to expire paste (default "12h")
  -u string
    	URL to post requests to (default "http://localhost:8080")
$ cat somefile | gp
http://localhost:8080/p/db14728a6185ca489095b791b1f40faa0d932ab5
$ gp -u https://paste.as65342.net somefile
https://paste.as65342.net/p/5cee9bc0c93784a8de492818da3123d4995980f9
$ export GP_URL="https://paste.as65342.net"
$ export GP_EXPIRY="1m"
$ cat anotherfile | gp
https://paste.as65342.net/p/cb99417c67825314f9ebd51f308286235c13142b
```

## Demo
A live demo of this software, running in a rkt container, can be seen on the URL below:

[paste.as65342.net](https://paste.as65342.net)

## Credits
This pastebin is heavily inspired by [maze.io/paste](https://maze.io/paste), hosted by [tehmaze](https://github.com/tehmaze).
