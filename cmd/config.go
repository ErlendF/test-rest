package cmd

import (
	"os"
	"test/pkg/database"

	"github.com/sirupsen/logrus"
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

const (
	configSetting   = "config"
	port            = "port"
	verbose         = "verbose"
	jsonFormatter   = "jsonFormatter"
	shutdownTimeout = "shutdownTimeout"
	dbType          = "dbType"
	dbName          = "dbName"
	dbUser          = "dbUser"
	dbPassword      = "dbPassword"
	dbSslMode       = "dbSslMode"
	dbHost          = "dbHost"
	dbPort          = "dbPort"
)

func init() {
	// Reads commandline arguments into config
	// Defaults are set in the setDefaults function
	pflag.StringP(configSetting, "c", "", "Config file")
	pflag.IntP(port, "p", 0, "Sets which port the application should listen to")
	pflag.BoolP(verbose, "v", false, "Verbose logging")
	pflag.BoolP(jsonFormatter, "j", false, "JSON logging format")
	pflag.IntP(shutdownTimeout, "s", 0,
		"Sets the timeout (in seconds) for graceful shutdown")

	// Database config
	pflag.String(dbType, "", "Database type (mysql, postgres)")
	pflag.String(dbName, "", "Name of the database the app should connect to")
	pflag.String(dbUser, "", "Username the app should use to connect to the database")
	pflag.String(dbPassword, "", "Password for the database user")
	pflag.String(dbSslMode, "",
		"SSL_MODE to use when connecting to the database (only used for PostgreSQL)")
	pflag.String(dbHost, "", "Host running the database")
	pflag.Int(dbPort, 0, "Port for connecting to the database")
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

	readConfigFile()

	// TODO: use viper to unmarshal config
	return config{
		verbose:         viper.GetBool(verbose),
		jsonFormatter:   viper.GetBool(jsonFormatter),
		shutdownTimeout: viper.GetInt(shutdownTimeout),
		port:            viper.GetInt(port),
		db: database.Config{
			Type:     viper.GetString(dbType),
			User:     viper.GetString(dbUser),
			Password: viper.GetString(dbPassword),
			SSLMode:  viper.GetString(dbSslMode),
			Host:     viper.GetString(dbHost),
			Name:     viper.GetString(dbName),
			Port:     viper.GetInt(dbPort),
		},
	}
}

func setDefaults() {
	viper.SetDefault(verbose, false)
	viper.SetDefault(jsonFormatter, false)
	viper.SetDefault(shutdownTimeout, 15)
	viper.SetDefault(port, 8080)
	viper.SetDefault(dbType, "postgres")
	viper.SetDefault(dbUser, "default")
	viper.SetDefault(dbPassword, "default")
	viper.SetDefault(dbSslMode, "disable")
	viper.SetDefault(dbHost, "localhost")
	viper.SetDefault(dbName, "postgres")
	viper.SetDefault(dbPort, 5432)
}

func bindEnvs() error {
	envs := [][]string{
		{configSetting},
		{verbose},
		{jsonFormatter, "JSON_FORMATTER"},
		{shutdownTimeout, "SHUTDOWN_TIMEOUT"},
		{port},
		{dbType, "DB_TYPE"},
		{dbUser, "DB_USER"},
		{dbPassword, "DB_PASSWORD"},
		{dbSslMode, "DB_SSLMODE"},
		{dbHost, "DB_HOST"},
		{dbName, "DB_NAME"},
		{dbPort, "DB_PORT"},
	}

	for _, e := range envs {
		err := viper.BindEnv(e...)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO: change format of file to make db and object
func readConfigFile() {
	fileProvided := true
	cfgfile := viper.GetString(configSetting)
	if cfgfile == "" {
		cfgfile = ".config.yml"
		fileProvided = false
	}
	viper.SetConfigFile(cfgfile)

	err := viper.ReadInConfig()
	if err == nil {
		logrus.Infof("Using config file: %s", viper.ConfigFileUsed())
		return
	}

	if !fileProvided {
		logrus.Infoln("No config file provided")
		return
	}

	logrus.WithError(err).Errorf("Error reading config file")
}
