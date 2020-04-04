/*
Copyright © 2019 Erlend Fonnes erlend.fonnes@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"test/pkg/database"
	"test/pkg/server"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type config struct {
	verbose         bool
	jsonFormatter   bool
	shutdownTimeout int
	port            int
	db              database.Config
}

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "test",
	Short: "Test-Rest is a REST API for a very simple message board.",
	Long: `Test-Rest is a REST API for a very simple message board.
It supports and uses a MySQL or PostgreSQL database for persistent storage.
See README.md for additional documentation.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := initConfig()
		setupLog(cfg.verbose, cfg.jsonFormatter)
		logrus.Debugf("Startup config: %+v", cfg)

		// Getting a database instance
		db, err := database.New(&cfg.db)
		if err != nil {
			logrus.WithError(err).Fatalf("Unable to get new Database:%s", err)
		}

		err = db.Migrate()
		if err != nil {
			logrus.Fatal(err)
		}

		srv := server.New(db, cfg.port)

		// Making an channel to listen for errors (later blocking until either error or signal is received)
		errChan := make(chan error)

		// Starting server in a go routine to allow for graceful shutdown and potentially additional services
		go func() {
			logrus.Infof("Starting server on port %d", cfg.port)
			if err := srv.ListenAndServe(); err != nil {
				errChan <- err
			}
		}()

		// Attempting to catch quit via SIGINT (Ctrl+C) to shut down gracefully
		// SIGKILL, SIGQUIT or SIGTERM will not be caught.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		// Blocking until signal or error is received
		select {
		case <-c:
			logrus.Infof("Shutting down server due to interrupt")
		case err := <-errChan:
			logrus.WithError(err).Errorf("Shutting down server due to error")
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.shutdownTimeout)*time.Second)
		defer cancel()

		// Attempting to shut down the server
		if err := srv.Shutdown(ctx); err != nil {
			logrus.WithError(err).Fatalf("Unable to gracefully shutdown server")
		}

		logrus.Infoln("Finished shutting down")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Reads commandline arguments into config
	// Defaults are set in the setDefaults function
	pflag.StringP("config", "c", "", "Config file")
	pflag.IntP("port", "p", 0, "Sets which port the application should listen to")
	pflag.BoolP("verbose", "v", false, "Verbose logging")
	pflag.BoolP("jsonFormatter", "j", false, "JSON logging format")
	pflag.IntP("shutdownTimeout", "s", 0,
		"Sets the timeout (in seconds) for graceful shutdown")

	// Database config
	pflag.String("dbtype", "", "Database type (mysql, postgres)")
	pflag.String("dbname", "", "Name of the database the app should connect to")
	pflag.String("dbuser", "", "Username the app should use to connect to the database")
	pflag.String("dbpassword", "", "Password for the database user")
	pflag.String("dbsslmode", "",
		"SSL_MODE to use when connecting to the database (only used for PostgreSQL)")
	pflag.String("dbhost", "", "Host running the database")
	pflag.Int("dbport", 0, "Port for connecting to the database")
}

// setupLog initializes logrus logger
func setupLog(verbose, jsonFormatter bool) {
	logLevel := logrus.InfoLevel

	if verbose {
		logLevel = logrus.DebugLevel
	}

	logrus.SetLevel(logLevel)
	logrus.SetOutput(os.Stdout)

	if jsonFormatter {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}

func initConfig() config {
	setDefaults()
	if err := bindEnvs(); err != nil {
		logrus.Fatalf("could not bind environment variables with spf13/viper: %v", err)
	}

	pflag.CommandLine.VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			if err := viper.BindPFlag(f.Name, f); err != nil {
				logrus.Fatalf("could not bind flag with spf13/viper: %v", err)
			}
		}
	})

	cfgfile := pflag.Lookup("config").Value.String()
	if cfgfile == "" {
		cfgfile = ".config.yml"
	}
	viper.SetConfigFile(cfgfile)

	if err := viper.ReadInConfig(); err == nil {
		logrus.Infof("Using config file: %s", viper.ConfigFileUsed())
	}

	return config{
		verbose:         viper.GetBool("verbose"),
		jsonFormatter:   viper.GetBool("jsonFormater"),
		shutdownTimeout: viper.GetInt("shutdownTimeout"),
		port:            viper.GetInt("port"),
		db: database.Config{
			Type:     viper.GetString("dbtype"),
			User:     viper.GetString("dbuser"),
			Password: viper.GetString("dbpassword"),
			SSLMode:  viper.GetString("dbsslmode"),
			Host:     viper.GetString("dbhost"),
			Name:     viper.GetString("dbname"),
			Port:     viper.GetInt("dbport"),
		},
	}
}

func setDefaults() {
	viper.SetDefault("verbose", false)
	viper.SetDefault("jsonFormater", false)
	viper.SetDefault("shutdownTimeout", 15)
	viper.SetDefault("port", 8080)
	viper.SetDefault("dbtype", "postgres")
	viper.SetDefault("dbuser", "default")
	viper.SetDefault("dbpassword", "default")
	viper.SetDefault("dbsslmode", "disable")
	viper.SetDefault("dbhost", "localhost")
	viper.SetDefault("dbname", "default")
	viper.SetDefault("dbport", 5432)
}

func bindEnvs() error {
	envs := [][]string{
		{"verbose"},
		{"jsonFormater", "JSON_FORMATER"},
		{"shutdownTimeout", "SHUTDOWN_TIMEOUT"},
		{"port"},
		{"dbtype", "DB_TYPE"},
		{"dbuser", "DB_USER"},
		{"dbpassword", "DB_PASSWORD"},
		{"dbsslmode", "DB_SSLMODE"},
		{"dbhost", "DB_HOST"},
		{"dbname", "DB_NAME"},
		{"dbport", "DB_PORT"},
	}

	for _, e := range envs {
		err := viper.BindEnv(e...)
		if err != nil {
			return err
		}
	}

	return nil
}
