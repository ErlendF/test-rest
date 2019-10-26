package database

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //Postgres driver
	"github.com/sirupsen/logrus"
)

//Database contains a sql db
type Database struct {
	*sqlx.DB
	UsernameNumberMax int
}

//New returns a database object
func New() (*Database, error) {
	db := &Database{UsernameNumberMax: 9999}
	var err error
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s",
		host, port, user, password, dbname)

	logrus.Debugf("Trying to connect to db: %s", psqlInfo)
	db.DB, err = sqlx.Connect("postgres", psqlInfo)
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
