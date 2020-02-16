package database

import (
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
)

func (db *Database) Migrate() error {
	migrate.SetTable("migrations")
	num, err := migrate.Exec(db.DB.DB, db.dbType, db.getMigrations(), migrate.Up)
	if err != nil {
		return fmt.Errorf("could not perform database migrations: %w", err)
	}

	logrus.Infof("Executed %d migrations", num)
	return nil
}
