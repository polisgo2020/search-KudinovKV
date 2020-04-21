package database

import (
	"github.com/go-pg/pg/v9"
	zl "github.com/rs/zerolog/log"
	"log"
)

type Index struct {
	FileName string `pg:"fileName,pk"`
	Token    string `pg:"token,pk"`
}

func InitDB(toConnect string, createOrNot bool) (*pg.DB, error) {
	pgOpt, err := pg.ParseURL(toConnect)
	if err != nil {
		return nil, err
	}
	pgdb := pg.Connect(pgOpt)
	if createOrNot == true {
		_, err = pgdb.Exec(`
		CREATE TABLE public.Index 
		(
			fileName text, 
			token text,
			PRIMARY KEY (fileName, token)
		);`)
		if err != nil {
			zl.Fatal().Err(err).
				Msg("Can't create database")
		}
	}
	return pgdb, nil
}

func AddIndex(in []Index, pg *pg.DB) error {
	/*var err error
	for _, element := range i {
		_, err = pg.Exec("INSERT INTO Index (fileName, token) VALUES (?, ?)", element.FileName, element.Token)
		if err != nil {
			return err
		}
	}*/
	_, err := pg.Model(&in).Insert()
	zl.Debug().
		Msgf("Insert %d elements ", len(in))
	return err
}

func GetFiles(token string, pg *pg.DB) ([]Index, error) {
	// Тут аналогично
	// "ERROR #42P01 relation \"indices\" does not exist"
	// err := m.pg.Model(&i).Where("token=?", token).Select()
	i := []Index{}
	err := pg.Model(&i).Where("token=?", token).Select()
	log.Println(token, i)
	return i, err
}
