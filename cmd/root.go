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

	_ "github.com/joho/godotenv/autoload" // importing .env to os.env
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var config struct {
	verbose         bool
	jsonFormatter   bool
	shutdownTimeout int
	version         int
	port            int
	dbType          string
}

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "test",
	Short: "Test",
	Long:  `Test`,
	Run: func(cmd *cobra.Command, args []string) {
		setupLog(config.verbose, config.jsonFormatter)
		logrus.Debugf("Startup config: %+v", config)

		// Getting a database instance
		db, err := database.New(config.dbType)
		if err != nil {
			logrus.WithError(err).Fatalf("Unable to get new Database:%s", err)
		}

		err = db.Migrate()
		if err != nil {
			logrus.Fatal(err)
		}

		srv := server.New(db, config.port)

		// Making an channel to listen for errors (later blocking until either error or signal is received)
		errChan := make(chan error)

		// Starting server in a go routine to allow for graceful shutdown and potentially additional services
		go func() {
			logrus.Infof("Starting server on port %d", config.port)
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

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.shutdownTimeout)*time.Second)
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
	rootCmd.Flags().IntVarP(&config.port, "port", "p", 8080, "Sets which port the application should listen to")
	rootCmd.Flags().IntVarP(&config.shutdownTimeout, "shutdownTimeout", "s", 15, "Sets the timeout (in seconds) for graceful shutdown")
	rootCmd.Flags().BoolVarP(&config.verbose, "verbose", "v", false, "Verbose logging")
	rootCmd.Flags().BoolVarP(&config.jsonFormatter, "jsonFormatter", "j", false, "JSON logging format")
	rootCmd.Flags().StringVarP(&config.dbType, "database", "d", "mysql", "Database type (mysql, postgres)")
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
