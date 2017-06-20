package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/r3boot/go-paste/lib"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"time"
)

const (
	D_DEBUG        bool   = false
	D_URL          string = "http://localhost:8080"
	D_EXPIRY       string = "1h"
	D_OPEN_URL     bool   = false
	PASTE_MAX_SIZE int64  = 1073741824
	PASTE_PREFIX   string = "/p/"
)

var (
	debug        = flag.Bool("D", D_DEBUG, "Enable debug output")
	cli_pastebin = flag.String("u", D_URL, "URL to post requests to")
	cli_expiry   = flag.String("e", D_EXPIRY, "When to expire paste")
	cli_open_url = flag.Bool("o", D_OPEN_URL, "Open the new url in a browser")
	pastebin     string
	expiry       string
	open_url     bool
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

	open_url = *cli_open_url
	if open_url == D_OPEN_URL {
		if value = os.Getenv("GP_OPEN_URL"); value != "" {
			switch value {
			case "y":
				{
					open_url = true
				}
			case "Y":
				{
					open_url = true
				}
			case "yes":
				{
					open_url = true
				}
			case "1":
				{
					open_url = true
				}
			default:
				{
					open_url = false
				}
			}
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
		err         error
		client      http.Client
		values      url.Values
		resp        *http.Response
		duration    time.Duration
		hash        string
		responseURL string
		cmd         *exec.Cmd
		opener      string
		output      []byte
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
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	values = url.Values{}
	values.Add("content", string(content))
	values.Add("expire", expiry)

	if resp, err = client.PostForm(pastebin, values); err != nil {
		Log.Error("Failed to post form: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusMovedPermanently {
		hash = resp.Header["Location"][0][len(PASTE_PREFIX):]
		responseURL = pastebin + "/p/" + hash
		fmt.Printf(responseURL + "\n")

		if open_url {
			switch runtime.GOOS {
			case "linux":
				{
					opener = "xdg-open"
				}
			case "darwin":
				{
					opener = "open"
				}
			default:
				{
					Log.Error("Open functionality not supported on " + runtime.GOOS)
				}
			}

			cmd = exec.Command(opener, responseURL)
			if output, err = cmd.CombinedOutput(); err != nil {
				Log.Warn(string(output))
				Log.Error("Failed to run " + opener + ": " + err.Error())
			}
		}
	} else {
		Log.Error("Received an invalid return code from pastebin: " + resp.Status)
	}

}
