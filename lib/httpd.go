package lib

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Default expiry options
var DEFAULT_EXPIRY_OPTIONS []string = []string{
	"1m",
	"15m",
	"30m",
	"1h",
	"1d",
	"1w",
	"2w",
	"4w",
	"3M",
	"6M",
	"1Y",
}

const DEFAULT_EXPIRY_OPTION string = "1d"

const PASTE_PREFIX = "/p/"

func viewPaste(w http.ResponseWriter, r *http.Request) {
	var err error
	var hash string
	var paste *Paste

	hash = r.URL.Path[len(PASTE_PREFIX):]

	if hash == "" {
		io.WriteString(w, "No such hash")
		return
	}

	if paste, err = LoadPaste(hash); err != nil {
		io.WriteString(w, "No such hash")
		return
	}

	w.Write(paste.Content)
}

func newPaste(w http.ResponseWriter, r *http.Request) {
	var fs os.FileInfo
	var err error
	var t *template.Template
	var p *NewPaste
	var paste *Paste
	var hash_path string
	var duration time.Duration

	if r.Method == "GET" {
		if fs, err = os.Stat(Config.Template); err != nil {
			Log.Warning("Cannot find template")
			return
		}

		if fs.IsDir() {
			Log.Warning(Config.Template + " is a directory")
			return
		}

		if t, err = template.ParseFiles(Config.Template); err != nil {
			Log.Warning("Failed to load template: " + err.Error())
			return
		}

		p = &NewPaste{
			Title:               Config.Title,
			ExpiryOptions:       DEFAULT_EXPIRY_OPTIONS,
			DefaultExpiryOption: DEFAULT_EXPIRY_OPTION,
		}

		t.Execute(w, p)
		Log.Debug("Served page on " + r.URL.Path)
	} else if r.Method == "POST" {
		duration, err = time.ParseDuration(r.PostFormValue("expire"))
		if err != nil {
			Log.Warning("Failed to parse duration")
			return
		}

		paste = &Paste{
			Content:    []byte(r.PostFormValue("content")),
			Expiration: duration,
		}

		paste.Save()

		hash_path = "/p/" + paste.Hash
		http.Redirect(w, r, hash_path, 301)
		Log.Debug("Redirected to " + hash_path)
	}
}

func serveAsset(w http.ResponseWriter, r *http.Request) {
	Log.Debug(r.URL.Path)
}

func SetupServer() {
	http.HandleFunc("/p/", viewPaste)
	http.HandleFunc("/", newPaste)
}

func RunServer() {
	var addr string

	addr = Config.BindIp + ":" + strconv.Itoa(Config.BindPort)

	Log.Debug("Listening on " + addr)
	http.ListenAndServe(addr, nil)
}
