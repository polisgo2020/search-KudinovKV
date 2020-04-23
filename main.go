package main

import (
	"github.com/go-pg/pg/v9"
	"github.com/polisgo2020/search-KudinovKV/database"
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

var (
	exampleBuild = "Invalid number of arguments. Example of call: <build> /path/to/files /path/to/output or" +
		" <build> database "
	exampleSearch = "Invalid number of arguments. Example of call: <search> /path/to/index/file or" +
		" <search> database"
)

func searchMain(cfg config.Config, filename string) {
	data, err := file.ReadFile(filename)
	if err != nil {
		zl.Fatal().Err(err).
			Msg("Can't read index file ")
	}
	web.StartServer(cfg, data, nil)
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

func buildMainDB(dirName string, pgdb *pg.DB) {
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
	i.WriteResultDB(pgdb)
}

func main() {
	cfg := config.LoadConfig()
	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		zl.Fatal().Err(err).Msgf("Can't parse loglevel")
	}
	zerolog.SetGlobalLevel(logLevel)
	zl.Logger = zl.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	zl.Debug().
		Msgf("%v", os.Args)

	if len(os.Args) < 2 {
		zl.Fatal().
			Msg(exampleBuild)
	}
	if os.Args[1] == "build" {
		if os.Args[2] == "file" {
			if len(os.Args) < 5 {
				zl.Fatal().
					Msg(exampleBuild)
			}
			buildMain(os.Args[3], os.Args[4])
		} else if os.Args[2] == "database" {
			pgdb, err := database.InitDB(cfg.PgSQL, true)
			if err != nil {
				zl.Fatal().
					Msg("Can't connect to database")
			}
			defer pgdb.Close()
			buildMainDB(os.Args[3], pgdb)
		}
	} else if os.Args[1] == "search" {
		if len(os.Args) < 3 {
			zl.Fatal().
				Msg(exampleSearch)
		}
		if os.Args[2] == "database" {
			pgdb, err := database.InitDB(cfg.PgSQL, false)
			if err != nil {
				zl.Fatal().
					Msg("Can't connect to database")
			}
			defer pgdb.Close()
			web.StartServer(cfg, "", pgdb)
		} else {
			searchMain(cfg, os.Args[2])
		}
	} else {
		zl.Fatal().
			Msg(exampleBuild)
	}
}
