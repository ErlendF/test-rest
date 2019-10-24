package dbmigrate

import (
	"os"
	"time"

	db "test/pkg/database"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
)

type statusRow struct {
	ID        string
	Migrated  bool
	AppliedAt time.Time
}

// DoMigrate do database migrations
func DoMigrate(url string, sqlDir string) error {
	logrus.Infof("Do migration for %s", url)

	dbx := db.New(url)

	m, err := Migrations(sqlDir)
	if err != nil {
		return err
	}

	ms := &migrate.MemoryMigrationSource{Migrations: m}

	migrate.SetTable("migrations")
	numMigrationsPerformed, err := migrate.Exec(dbx.DB.DB, "postgres", ms, migrate.Up)
	if err != nil {
		return errors.Wrap(err, "could not perform database migrations")
	}

	logrus.Infof("Executed %d migrations", numMigrationsPerformed)

	return nil
}

// ShowMigrations shows executed migrations
func ShowMigrations(url string, sqlDir string) error {

	dbx := db.New(url)

	records, err := migrate.GetMigrationRecords(dbx.DB.DB, "postgres")
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Migration", "Applied"})
	table.SetColWidth(60)

	rows := make(map[string]*statusRow)
	ms, err := Migrations(sqlDir)
	if err != nil {
		return err
	}

	for _, m := range ms {
		rows[m.Id] = &statusRow{
			ID:       m.Id,
			Migrated: false,
		}
	}

	for _, r := range records {
		if rows[r.Id] == nil {
			logrus.Warnf("Could not find migration file: %v", r.Id)
			continue
		}

		rows[r.Id].Migrated = true
		rows[r.Id].AppliedAt = r.AppliedAt
	}

	for _, m := range ms {
		if rows[m.Id] != nil && rows[m.Id].Migrated {
			table.Append([]string{
				m.Id,
				rows[m.Id].AppliedAt.String(),
			})
		} else {
			table.Append([]string{
				m.Id,
				"no",
			})
		}
	}

	table.Render()
	return nil
}
