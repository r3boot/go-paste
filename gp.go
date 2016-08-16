package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/r3boot/rlib/logger"
	"io"
	"net/http"
	"net/url"
	"os"
)

const D_DEBUG bool = false
const D_URL string = "http://localhost:8080"
const D_EXPIRY string = "12h"

const PASTE_MAX_SIZE int64 = 1073741824
const PASTE_PREFIX string = "/p/"

var debug = flag.Bool("D", D_DEBUG, "Enable debug output")
var cli_pastebin = flag.String("u", D_URL, "URL to post requests to")
var cli_expiry = flag.String("e", D_EXPIRY, "When to expire paste")

var pastebin string
var expiry string

var Log logger.Log

var content []byte

func readFile(fname string) (content []byte, err error) {
	var fs os.FileInfo
	var fd *os.File
	var bytes_read int

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
	var reader *bufio.Reader
	var num_read int

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
	var err error
	var fi os.FileInfo
	var value string

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
		Log.Fatal("main: os.Stdin.Stat() failed: " + err.Error())
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		if len(flag.Args()) == 0 {
			Log.Fatal("init: Need something to do!")
		}
		content, err = readFile(flag.Arg(0))
	} else {
		content, err = readStdin()
	}
	if err != nil {
		Log.Fatal(err)
	}
}

func main() {
	var err error
	var client http.Client
	var values url.Values
	var resp *http.Response
	var hash string

	client = http.Client{}

	values = url.Values{}
	values.Add("content", string(content))
	values.Add("expire", expiry)

	if resp, err = client.PostForm(pastebin, values); err != nil {
		Log.Fatal("main -> client.PostForm() failed: " + err.Error())
	}

	if resp.StatusCode == 301 {
		hash = resp.Header["Location"][0][len(PASTE_PREFIX):]
		fmt.Printf(pastebin + "/p/" + hash + "\n")
	} else {
		Log.Fatal("main -> resp.StatusCode != 301: " + resp.Status)
	}
}
