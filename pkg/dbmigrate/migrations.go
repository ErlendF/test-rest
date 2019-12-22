package dbmigrate

import (
	"os"
	"path"
	"sort"
	"strings"

	migrate "github.com/rubenv/sql-migrate"
)

// Migrations sets up tables initially. Loads migrations from a directory
func Migrations(srcDir string) ([]*migrate.Migration, error) {
	migrations := make([]*migrate.Migration, 0)

	file, err := os.Open(srcDir)
	if err != nil {
		return nil, err
	}

	files, err := file.Readdir(0)
	if err != nil {
		return nil, err
	}

	for _, info := range files {
		if strings.HasPrefix(info.Name(), "mi_") {
			var file *os.File
			var migration *migrate.Migration
			if strings.HasSuffix(info.Name(), ".sql") {
				file, err = os.Open(path.Join(srcDir, info.Name()))
				if err != nil {
					return nil, err
				}
				migration, err = migrate.ParseMigration(info.Name(), file)
				if err != nil {
					return nil, err
				}
			} else {
				continue
			}

			migrations = append(migrations, migration)
		}
	}

	// Make sure migrations are sorted
	sort.Sort(byID(migrations))

	return migrations, nil
}

type byID []*migrate.Migration

func (b byID) Len() int           { return len(b) }
func (b byID) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byID) Less(i, j int) bool { return b[i].Less(b[j]) }
