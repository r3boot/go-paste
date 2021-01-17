package lib

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
)

// Default expiry options
var defaultExpiryOptions map[string]string = map[string]string{
	"5m":   "5 minutes",
	"1h":   "1 hour",
	"4h":   "4 hours",
	"24h":  "1 day",
	"168h": "1 week",
	"720h": "30 days",
}

const (
	defaultExpiryOption string = "1d"
	pastePrefix                = "/p/"
	clfTimeFormat       string = "02/Jan/2006:15:04:05 -0700"
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
	var (
		err    error
		errmsg string
		hash   string
		paste  *Paste
	)

	hash = r.URL.Path[len(pastePrefix):]

	if hash == "" {
		errmsg = "No hash specified"
		http.Error(w, errmsg, http.StatusNotFound)
		return
	}

	if paste, err = LoadPaste(hash); err != nil {
		errmsg = "No such hash"
		http.Error(w, errmsg, http.StatusNotFound)
		log.Error(errmsg + ": " + err.Error())
		return
	}

	w.Write(paste.Content)
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
			return
		}

		if t, err = template.New("index").Parse(templateData); err != nil {
			errmsg = "Failed to parse template"
			http.Error(w, errmsg, http.StatusInternalServerError)
			log.Error(errmsg + ": " + err.Error())
		}

		p = &NewPaste{
			Title:         Config.Title,
			ExpiryOptions: defaultExpiryOptions,
		}

		if err = t.Execute(&content, p); err != nil {
			errmsg = "Failed to render template"
			http.Error(w, errmsg, http.StatusInternalServerError)
			log.Error(errmsg + ": " + err.Error())
		}

		w.Write(content.Bytes())

	} else if r.Method == "POST" {
		duration, err = time.ParseDuration(r.PostFormValue("expire"))
		if err != nil {
			errmsg = "Failed to parse duration"
			http.Error(w, errmsg, http.StatusInternalServerError)
			log.Warn(errmsg + ": " + err.Error())
			return
		}

		paste = &Paste{
			Content:    []byte(r.PostFormValue("content")),
			Expiration: duration,
		}

		_, err = paste.Save()
		if err != nil {
			errmsg = "Failed to save paste"
			http.Error(w, errmsg, http.StatusInternalServerError)
			log.Warn(errmsg + ": " + err.Error())
			return
		}

		hash_path = "/p/" + paste.Hash
		http.Redirect(w, r, hash_path, http.StatusMovedPermanently)
	}
}

func RunServer() {
	var (
		addr string
	)

	addr = Config.BindIp + ":" + strconv.Itoa(Config.BindPort)

	http.HandleFunc("/p/", viewPaste)
	http.HandleFunc("/", newPaste)

	log.Debug("Listening on " + addr)
	http.ListenAndServe(addr, handlers.CombinedLoggingHandler(
		os.Stdout, http.DefaultServeMux))
}
