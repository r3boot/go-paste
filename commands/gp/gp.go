package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/r3boot/go-paste/lib"
	"io"
	"time"
	"net/http"
	"net/url"
	"os"
)

const (
	D_DEBUG        bool   = false
	D_URL          string = "http://localhost:8080"
	D_EXPIRY       string = "12h"
	PASTE_MAX_SIZE int64  = 1073741824
	PASTE_PREFIX   string = "/p/"
)

var (
	debug        = flag.Bool("D", D_DEBUG, "Enable debug output")
	cli_pastebin = flag.String("u", D_URL, "URL to post requests to")
	cli_expiry   = flag.String("e", D_EXPIRY, "When to expire paste")
	pastebin     string
	expiry       string
	Log          lib.Log
	content      []byte
)

func readFile(fname string) (content []byte, err error) {
	var (
		fs         os.FileInfo
		fd         *os.File
		bytes_read int
	)

	if fs, err = os.Stat(fname); err != nil {
		err = errors.New("readFile -> os.Stat(fname) failed: " + err.Error())
		return
	}
	if fs.IsDir() {
		err = errors.New("readFile -> fs.IsDir(): Is a directory")
		return
	}

	if fd, err = os.Open(fname); err != nil {
		err = errors.New("readFile -> os.Open(fname) failed: " + err.Error())
		return
	}

	content = make([]byte, fs.Size())
	if bytes_read, err = fd.Read(content); err != nil {
		content = nil
		err = errors.New("readfile -> fd.Read(content) failed: " + err.Error())
		return
	}

	if int64(bytes_read) != fs.Size() {
		content = nil
		err = errors.New("readFile -> bytes_read != fs.Size(): Corrupt read")
		return
	}

	return
}

func readStdin() (content []byte, err error) {
	var (
		reader   *bufio.Reader
		num_read int
	)

	reader = bufio.NewReader(os.Stdin)
	content = make([]byte, 0, PASTE_MAX_SIZE)

	num_read, err = reader.Read(content[:cap(content)])
	content = content[:num_read]
	if num_read == 0 {
		if err != nil && err != io.EOF {
			err = errors.New("readStdin -> reader.Read: 0 bytes read: " + err.Error())
			return
		}
	}

	return
}

func init() {
	var (
		err   error
		fi    os.FileInfo
		value string
	)

	flag.Parse()

	pastebin = *cli_pastebin
	if pastebin == D_URL {
		if value = os.Getenv("GP_URL"); value != "" {
			pastebin = value
		}
	}

	expiry = *cli_expiry
	if expiry == D_EXPIRY {
		if value = os.Getenv("GP_EXPIRY"); value != "" {
			expiry = value
		}
	}

	Log.UseDebug = *debug
	Log.UseVerbose = *debug
	Log.UseTimestamp = false
	Log.Debug("Logging initialized")

	fi, err = os.Stdin.Stat()
	if err != nil {
		Log.Error("main: os.Stdin.Stat() failed: " + err.Error())
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		if len(flag.Args()) == 0 {
			Log.Error("init: Need something to do!")
		}
		content, err = readFile(flag.Arg(0))
	} else {
		content, err = readStdin()
	}
	if err != nil {
		Log.Error(err)
	}
}

func main() {
	var (
		err    error
		client http.Client
		values url.Values
		resp   *http.Response
		duration time.Duration
		hash   string
	)

	if duration, err = time.ParseDuration(expiry); err != nil {
		Log.Error("Failed to parse duration: " + err.Error())
	}

	if duration < lib.EXPIRE_MIN {
		Log.Error("Duration needs to be larger then 1 minute")
	}

	if duration > lib.EXPIRE_MAX {
		Log.Error("Duration needs to be smaller then 60 days")
	}

	client = http.Client{}

	values = url.Values{}
	values.Add("content", string(content))
	values.Add("expire", expiry)

	if resp, err = client.PostForm(pastebin, values); err != nil {
		Log.Error("Failed to post form: " + err.Error())
	}

	if resp.StatusCode == 301 {
		hash = resp.Header["Location"][0][len(PASTE_PREFIX):]
		fmt.Printf(pastebin + "/p/" + hash + "\n")
	} else {
		Log.Error("Received a non 301 return code from pastebin: " + resp.Status)
	}
}
