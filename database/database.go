package database

import (
	"context"
	"github.com/go-pg/pg/v9"
	zl "github.com/rs/zerolog/log"
)

// Index is the container for an index in PgSQL.
type Index struct {
	tableName struct{} `pg:"Index"`
	FileName  string   `pg:"fileName,pk"`
	Token     string   `pg:"token,pk"`
}

type dbLogger struct{}

// BeforeQuery logging before query
func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

// AfterQuery logging after query
func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	uq, err := q.FormattedQuery()
	if err != nil {
		return err
	}
	zl.Debug().
		Str("query", uq).
		Msg("query")
	return nil
}

// InitDB parse env, create and connect to database
func InitDB(toConnect string, createOrNot bool) (*pg.DB, error) {
	pgOpt, err := pg.ParseURL(toConnect)
	if err != nil {
		return nil, err
	}
	pgdb := pg.Connect(pgOpt)
	if createOrNot == true {
		_, err = pgdb.Exec(` DROP TABLE public.Index;
		CREATE TABLE public.Index 
		(
			"fileName" text,
			"token" text,
			PRIMARY KEY ("fileName", "token")
		);`)
		if err != nil {
			zl.Fatal().Err(err).
				Msg("Can't create database")
		}
	}
	pgdb.AddQueryHook(dbLogger{})
	return pgdb, nil
}

// AddIndex add filename, token in database
func AddIndex(in []Index, pg *pg.DB) error {
	_, err := pg.Model(&in).Insert()
	zl.Debug().
		Msgf("Insert %d elements ", len(in))
	return err
}

// GetFiles get filenames from database
func GetFiles(token string, pg *pg.DB) ([]Index, error) {
	i := []Index{}
	err := pg.Model(&i).Where("token=?", token).Select()
	zl.Debug().
		Msgf("Get %d elements ", len(i))
	return i, err
}
