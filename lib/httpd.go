package lib

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Default expiry options
var DEFAULT_EXPIRY_OPTIONS map[string]string = map[string]string{
	"5m":   "5 minutes",
	"1h":   "1 hour",
	"4h":   "4 hours",
	"24h":  "1 day",
	"168h": "1 week",
	"720h": "30 days",
}

const (
	DEFAULT_EXPIRY_OPTION string = "1d"
	PASTE_PREFIX                 = "/p/"
	TF_CLF                string = "02/Jan/2006:15:04:05 -0700"
)

func httpLog(r *http.Request, code int, size int) {
	var (
		srcip   string
		logline string
	)

	srcip = r.Header.Get("X-Forwarded-For")
	if srcip == "" {
		srcip = r.RemoteAddr
	}

	logline = srcip + " - - [" + time.Now().Format(TF_CLF) + "] "
	logline += "\"" + r.Method + " " + r.URL.Path + " " + r.Proto + "\" "
	logline += strconv.Itoa(code) + " " + strconv.Itoa(size)

	fmt.Println(logline)
}

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
	var (
		err    error
		errmsg string
		hash   string
		paste  *Paste
	)

	hash = r.URL.Path[len(PASTE_PREFIX):]

	if hash == "" {
		errmsg = "No hash specified"
		http.Error(w, errmsg, http.StatusNotFound)
		httpLog(r, http.StatusNotFound, len(errmsg))
		return
	}

	if paste, err = LoadPaste(hash); err != nil {
		errmsg = "No such hash"
		http.Error(w, errmsg, http.StatusNotFound)
		httpLog(r, http.StatusNotFound, len(errmsg))
		return
	}

	w.Write(paste.Content)
	httpLog(r, http.StatusOK, len(paste.Content))
}

func newPaste(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		errmsg    string
		t         *template.Template
		p         *NewPaste
		paste     *Paste
		hash_path string
		duration  time.Duration
		srcip     string
		content   bytes.Buffer
	)

	srcip = r.Header.Get("X-Forwarded-For")
	if srcip == "" {
		srcip = r.RemoteAddr
	}

	if r.Method == "GET" {
		if r.URL.Path != "/" {
			errmsg = "No such file or directory"
			http.Error(w, errmsg, http.StatusNotFound)
			httpLog(r, http.StatusNotFound, len(errmsg))
			return
		}

		if t, err = template.New("index").Parse(TEMPLATE_DATA); err != nil {
			errmsg = "Failed to parse template"
			http.Error(w, errmsg, http.StatusInternalServerError)
			httpLog(r, http.StatusInternalServerError, len(errmsg))
			log.Error(errmsg + ": " + err.Error())
		}

		p = &NewPaste{
			Title:         Config.Title,
			ExpiryOptions: DEFAULT_EXPIRY_OPTIONS,
		}

		if err = t.Execute(&content, p); err != nil {
			errmsg = "Failed to render template"
			http.Error(w, errmsg, http.StatusInternalServerError)
			httpLog(r, http.StatusInternalServerError, len(errmsg))
			log.Error(errmsg + ": " + err.Error())
		}

		w.Write(content.Bytes())
		httpLog(r, http.StatusOK, content.Len())

	} else if r.Method == "POST" {
		duration, err = time.ParseDuration(r.PostFormValue("expire"))
		if err != nil {
			errmsg = "Failed to parse duration"
			http.Error(w, errmsg, http.StatusInternalServerError)
			httpLog(r, http.StatusInternalServerError, len(errmsg))
			log.Warn(errmsg + ": " + err.Error())
			return
		}

		paste = &Paste{
			Content:    []byte(r.PostFormValue("content")),
			Expiration: duration,
		}

		paste.Save()

		hash_path = "/p/" + paste.Hash
		http.Redirect(w, r, hash_path, http.StatusMovedPermanently)
		httpLog(r, http.StatusMovedPermanently, 0)
	}
}

func SetupServer() {
	http.HandleFunc("/p/", viewPaste)
	http.HandleFunc("/", newPaste)
}

func RunServer() {
	var (
		addr string
	)

	addr = Config.BindIp + ":" + strconv.Itoa(Config.BindPort)

	log.Debug("Listening on " + addr)
	http.ListenAndServe(addr, nil)
}
