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
	"context"
	"os"
	"os/signal"
	"test/pkg/server"
	"time"

	"github.com/sirupsen/logrus"

	_ "github.com/joho/godotenv/autoload" //importing .env to os.env
	"github.com/spf13/cobra"
)

var config struct {
	verbose         bool
	jsonFormatter   bool
	shutdownTimeout int
	version         int
	port            int
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the rsp server",
	Long:  `Starts the rsp server to serve REST API.`,
	Run: func(cmd *cobra.Command, args []string) {
		setupLog(config.verbose, config.jsonFormatter)
		logrus.Debugf("Startup config: %+v", config)

		srv := server.New(config.port)

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

func init() {
	rootCmd.AddCommand(serveCmd)

}
