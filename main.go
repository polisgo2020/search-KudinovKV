package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/polisgo2020/search-KudinovKV/config"
	"github.com/polisgo2020/search-KudinovKV/file"
	"github.com/polisgo2020/search-KudinovKV/index"
	"github.com/rs/zerolog"
	zl "github.com/rs/zerolog/log"
)

var (
	listOfFiles []string
	i           *index.InvertIndex
)

func searchPage(w http.ResponseWriter, r *http.Request) {
	tokens := r.FormValue("tokens")
	if tokens == "" {
		zl.Debug().
			Msg("Incorrect request, cant find tokens field")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	out := i.MakeSearch(index.PrepareTokens(tokens), listOfFiles)
	tmpl := template.Must(template.ParseFiles("template/search.html"))
	tmpl.Execute(w,
		struct {
			Result []index.Rate
		}{
			out,
		})
}

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

func indexPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/index.html"))
	tmpl.Execute(w, struct{}{})
}

func searchMain(cfg config.Config, filename string) {
	data, err := file.ReadFile(filename)
	if err != nil {
		zl.Fatal().Err(err).Msg("Can't read index file ")
	}

	i = index.NewInvertIndex()
	listOfFiles = i.ParseIndexFile(data)

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexPage)
	mux.HandleFunc("/search", searchPage)
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

func buildMain(dirName, resultFilename string) {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		zl.Fatal().Err(err).
			Msg("Can't open dir")
	}
	i := index.NewInvertIndex()
	for _, f := range files {
		i.GetWg().Add(1)
		go i.MakeBuild(filepath.Join(dirName, f.Name()))
	}
	i.GetWg().Wait()
	i.WriteResult(resultFilename)
}

func main() {
	cfg := config.LoadConfig()
	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		zl.Fatal().Err(err).Msgf("Can't parse loglevel")
	}
	zerolog.SetGlobalLevel(logLevel)
	zl.Logger = zl.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if len(os.Args) < 2 {
		zl.Fatal().
			Msg("Invalid number of arguments. Example of call: <build> /path/to/files /path/to/output or" +
				" <search> /path/to/index/file")
	}
	if os.Args[1] == "build" {
		if len(os.Args) < 4 {
			zl.Fatal().
				Msg("Invalid number of arguments. Example of call: <build> /path/to/files /path/to/output")
		}
		buildMain(os.Args[2], os.Args[3])
	} else if os.Args[1] == "search" {
		if len(os.Args) < 3 {
			zl.Fatal().
				Msg("Invalid number of arguments. Example of call: <search> /path/to/index/file")
		}
		searchMain(cfg, os.Args[2])
	} else {
		zl.Fatal().
			Msg("Invalid number of arguments. Example of call: <build> /path/to/files /path/to/output or " +
				" <search> /path/to/index/file")
	}
}
