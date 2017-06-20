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

const (
	DEFAULT_EXPIRY_OPTION string = "1d"
	PASTE_PREFIX                 = "/p/"
)

func response(w http.ResponseWriter, r *http.Request, msg string) {
	var srcip string

	srcip = r.Header.Get("X-Forwarded-For")
	if srcip == "" {
		srcip = r.RemoteAddr
	}

	log.Info("[" + srcip + "]: " + msg)
	io.WriteString(w, msg)
}

func viewPaste(w http.ResponseWriter, r *http.Request) {
	var err error
	var hash string
	var paste *Paste

	hash = r.URL.Path[len(PASTE_PREFIX):]

	if hash == "" {
		response(w, r, "No such hash")
		return
	}

	if paste, err = LoadPaste(hash); err != nil {
		response(w, r, "No such hash")
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
	var srcip string

	srcip = r.Header.Get("X-Forwarded-For")
	if srcip == "" {
		srcip = r.RemoteAddr
	}

	if r.Method == "GET" {
		if fs, err = os.Stat(Config.Template); err != nil {
			response(w, r, "Cannot find template: "+Config.Template)
			return
		}

		if fs.IsDir() {
			response(w, r, Config.Template+" is a directory")
			return
		}

		if t, err = template.ParseFiles(Config.Template); err != nil {
			response(w, r, "Failed to load template: "+err.Error())
			return
		}

		p = &NewPaste{
			Title:         Config.Title,
			ExpiryOptions: DEFAULT_EXPIRY_OPTIONS,
		}

		t.Execute(w, p)
		log.Info("[" + srcip + "]: 200 " + r.URL.Path)
	} else if r.Method == "POST" {
		duration, err = time.ParseDuration(r.PostFormValue("expire"))
		if err != nil {
			log.Warn("Failed to parse duration: " + r.PostFormValue("expire"))
			io.WriteString(w, "Failed to parse duration: "+r.PostFormValue("expire"))
			return
		}

		paste = &Paste{
			Content:    []byte(r.PostFormValue("content")),
			Expiration: duration,
		}

		paste.Save()

		hash_path = "/p/" + paste.Hash
		http.Redirect(w, r, hash_path, 301)
		log.Info("[" + srcip + "]: 301 " + hash_path)
	}
}

func serveAsset(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.URL.Path)
}

func SetupServer() {
	http.HandleFunc("/p/", viewPaste)
	http.HandleFunc("/", newPaste)
}

func RunServer() {
	var addr string

	addr = Config.BindIp + ":" + strconv.Itoa(Config.BindPort)

	log.Debug("Listening on " + addr)
	http.ListenAndServe(addr, nil)
}
