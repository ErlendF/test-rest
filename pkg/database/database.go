package database

import (
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/sirupsen/logrus"
)

// Database contains a sql db
type Database struct {
	*sqlx.DB
	dbType string
}

// New returns a database object
func New(dbType string) (*Database, error) {
	db := &Database{dbType: dbType}
	var err error
	var info string

	switch dbType {
	case "postgres":
		info = postgreSQLInfo()
	case "mysql":
		info = mySQLInfo()
	default:
		return nil, fmt.Errorf("not supported database type: %s", dbType)
	}

	logrus.Infof("Trying to connect to db: %s", info)
	db.DB, err = sqlx.Connect(dbType, info)
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

func mySQLInfo() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbname)
}

func postgreSQLInfo() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
}
