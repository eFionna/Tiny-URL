package app

import (
	"log"
	"net/http"
	"strings"

	"github.com/eFionna/Tiny-URL/internal/url"
)

type PageData struct {
	ShortURL string
	Error    string
}

func (a *App) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := PageData{}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			data.Error = "Invalid form submission"
		} else {
			rawURL := r.Form.Get("url")
			code, err := url.ShortenURL(a.Ctx, a.RDB, rawURL)
			if err != nil {
				data.Error = err.Error()
			} else {
				data.ShortURL = a.Config.BaseURL + "/s/" + code
			}
		}
	}

	if err := a.Tmpl.Execute(w, data); err != nil {
		log.Println("Template rendering error:", err)
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}

func (a *App) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/s/")
	if code == "" {
		http.NotFound(w, r)
		return
	}

	longURL, err := url.GetURL(a.Ctx, a.RDB, code)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}
