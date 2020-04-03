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

type Config struct {
	Type     string
	User     string
	Password string
	SSLMode  string
	Host     string
	Name     string
	Port     int
}

// New returns a database object
func New(cfg *Config) (*Database, error) {
	db := &Database{dbType: cfg.Type}
	var err error
	var info string

	switch cfg.Type {
	case "postgres":
		info = postgreSQLInfo(cfg)
	case "mysql":
		info = mySQLInfo(cfg)
	default:
		return nil, fmt.Errorf("not supported database type: %s", cfg.Type)
	}

	logrus.Infof("Trying to connect to db: %s", info)
	db.DB, err = sqlx.Connect(cfg.Type, info)
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

func postgreSQLInfo(cfg *Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
}

func mySQLInfo(cfg *Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
}
