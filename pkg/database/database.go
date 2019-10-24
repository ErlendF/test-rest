package database

import (
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

//Database contains a sql db
type Database struct {
	*sqlx.DB
	UsernameNumberMax int
}

//New returns a database object
func New(URL string) (*Database, error) {
	db := &Database{UsernameNumberMax: 9999}
	var err error
	db.DB, err = sqlx.Connect("postgres", URL)
	if err != nil {
		return nil, err
	}

	maxOpen := 20
	maxIdle := 15

	if max, err := strconv.Atoi(os.Getenv("DB_MAXOPEN")); err == nil {
		logrus.Infof("Setting max open connections to %d", max)
		maxOpen = max
	}

	if max, err := strconv.Atoi(os.Getenv("DB_MAXIDLE")); err == nil {
		logrus.Infof("Setting max idle connections to %d", max)
		maxIdle = max
	}

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)

	return db, nil
}
