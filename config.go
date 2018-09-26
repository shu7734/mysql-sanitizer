package main

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

const usageString = "Usage: mysql-sanitizer [-v log-level] [-o output] [-p local-port] config-file"

// Config collects all the daemon's configuration options.
type Config struct {
	LogFile       string // The logfile we're writing to
	MysqlHost     string // The host running MySQL
	MysqlPort     int    // The MySQL server port on the MySQL host
	MysqlUsername string // The username to log into MySQL with
	MysqlPassword string // The password to log into MySQL with
	ListeningPort int    // The port to listen for client connections on
	LogLevel      int    // How much output to generate
}

var defaultConfig = Config{
	"-",         // LogFile
	"localhost", // MysqlHost
	3306,        // MysqlPort
	"root",      // MysqlUsername
	"",          // MysqlPassword
	3306,        // ListeningPort
	0,           // LogLevel
}

// GetConfig returns a compendium of configurations collected from the command line.
func GetConfig() Config {
	config := defaultConfig
	var configFile string

	switch len(flag.Args()) {
	case 0:
		configFile = os.Getenv("HOME") + "/.mysql-sanitizer.conf"
	case 1:
		configFile = flag.Arg(0)
	default:
		log.Fatal(usageString)
	}
	verifyConfigPermissions(configFile)

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatalf("Couldn't read config file %s: %s", configFile, err)
	}

	if config.MysqlUsername == "" {
		log.Fatal("No MysqlUsername found in the config file!")
	}

	// Read the command-line flags.
	flag.StringVar(&config.LogFile, "o", "-", "The filename to log output to (default stdout)")
	flag.IntVar(&config.ListeningPort, "p", config.ListeningPort, "The port to listen for client connections on (default 3306)")
	flag.IntVar(&config.LogLevel, "v", config.LogLevel, "The verbosity level (0-3, default 0)")
	flag.Parse()

	return config
}

// Throws an error if the file is readable by anyone but the user.
func verifyConfigPermissions(configFile string) {
	info, err := os.Stat(configFile)
	if err != nil {
		log.Fatalf("Can't stat the config file %s: %s", configFile, err)
	}

	if info.Mode()&0077 > 0 {
		log.Fatalf("The config file has excessively permissive permissions! Try \"chmod 0600 %s\".", configFile)
	}
}
