package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/sirupsen/logrus"
)

// Database contains a sql db
type Database struct {
	*sqlx.DB
}

// New returns a database object
func New(dbType string) (*Database, error) {
	db := &Database{}
	var err error
	var dbOpen *sqlx.DB

	switch dbType {
	case "postgres":
		dbOpen, err = openPostgreSQL()
	case "mysql":
		dbOpen, err = openMySQL()
	default:
		return nil, fmt.Errorf("not supported database type: %s", dbType)
	}
	if err != nil {
		return nil, err
	}

	err = dbOpen.Ping()
	for i := 0; i < 30 && err != nil; i++ {
		time.Sleep(time.Duration(i*2) * time.Second)
		logrus.Warnf("could not ping database, retrying:")
		err = dbOpen.Ping()
	}
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %v", err)
	}

	db.DB = dbOpen
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

func openMySQL() (*sqlx.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

	logrus.Debugf("Trying to connect to db: %s", uri)
	return sqlx.Open("mysql", uri)
}

func openPostgreSQL() (*sqlx.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	logrus.Debugf("Trying to connect to db: %s", psqlInfo)
	return sqlx.Open("postgres", psqlInfo)
}
