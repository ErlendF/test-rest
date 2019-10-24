// Copyright © 2019 Erlend Fonnes erlend.fonnes@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"test/pkg/dbmigrate"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

type dbMigrateConfig struct {
	ShowOnly bool
	SQLDir   string
}

var dbmConfig dbMigrateConfig

// dbMigrateCmd represents the dbMigrate command
var dbMigrateCmd = &cobra.Command{
	Use:   "dbMigrate",
	Short: "Migrate database to last config",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dbURL := os.Getenv("DB_URL")
		if dbURL == "" {
			logrus.Fatal("No DB_URL specified, required")
		}
		setupLog(config.verbose, config.jsonFormatter)
		if dbmConfig.ShowOnly {
			err := dbmigrate.ShowMigrations(dbURL, dbmConfig.SQLDir)
			if err != nil {
				logrus.Warn(err)
			}
			return
		}

		err := dbmigrate.DoMigrate(dbURL, dbmConfig.SQLDir)
		if err != nil {
			logrus.Warn(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(dbMigrateCmd)
}
