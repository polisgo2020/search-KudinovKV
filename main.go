package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/polisgo2020/search-KudinovKV/config"
	"github.com/polisgo2020/search-KudinovKV/file"
	"github.com/polisgo2020/search-KudinovKV/index"
	"github.com/polisgo2020/search-KudinovKV/web"
	"github.com/rs/zerolog"
	zl "github.com/rs/zerolog/log"
)

func searchMain(cfg config.Config, filename string) {
	data, err := file.ReadFile(filename)
	if err != nil {
		zl.Fatal().Err(err).
			Msg("Can't read index file ")
	}

	web.StartServer(cfg, data)
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
