package web

import (
	"html/template"
	"net/http"
	"time"

	"github.com/polisgo2020/search-KudinovKV/config"
	"github.com/polisgo2020/search-KudinovKV/index"
	zl "github.com/rs/zerolog/log"
)

var (
	searchTmpl  *template.Template
	listOfFiles []string
	i           *index.InvertIndex
)

// StartServer parse html templates, set handle func,set log middleware, create index and started server
func StartServer(cfg config.Config, data string) {
	indexTmpl, err := template.ParseFiles("web/index.html")
	if err != nil {
		zl.Fatal().Err(err).
			Msg("Can't parse index template")
	}
	searchTmpl, err = template.ParseFiles("web/search.html")
	if err != nil {
		zl.Fatal().Err(err).
			Msg("Can't parse search template")
	}

	i = index.NewInvertIndex()
	listOfFiles = i.ParseIndexFile(data)

	mux := http.NewServeMux()
	mux.HandleFunc("/search", searchPage)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err = indexTmpl.ExecuteTemplate(w, "index.html", struct{}{})
		if err != nil {
			zl.Fatal().Err(err).
				Msg("Can't execute index template")
		}
	})

	siteHandler := logMiddleware(mux)
	server := http.Server{
		Addr:         cfg.Listen,
		Handler:      siteHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	zl.Debug().
		Msg("[ + ] starting server at " + cfg.Listen)

	server.ListenAndServe()
}

// searchPage handle request, find tokens and execute search template
func searchPage(w http.ResponseWriter, r *http.Request) {
	tokens := r.FormValue("tokens")
	if tokens == "" {
		zl.Debug().
			Msg("Incorrect request, cant find tokens field")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	out := i.MakeSearch(index.PrepareTokens(tokens), listOfFiles)

	err := searchTmpl.ExecuteTemplate(w, "search.html",
		struct {
			Result []index.Rate
		}{
			out,
		})

	if err != nil {
		zl.Fatal().Err(err).
			Msg("Can't execute search template")
	}
}

// logMiddleware logging all request
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		zl.Debug().
			Str("method", r.Method).
			Str("remote", r.RemoteAddr).
			Str("path", r.URL.Path).
			Int("duration", int(time.Since(start))).
			Msgf("Called url %s", r.URL.Path)
	})
}
