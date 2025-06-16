package internal

import (
	"cmp"
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Host      string
	Port      int
	Debug     bool
	DBConnStr string
}

const (
	defaultHost = "0.0.0.0"
	defaultPort = 8080
	defaultDB   = "postgres://user.password@localhost:5432/notes?sslmode=disable"
)

func ReadConfig() (*Config, error) {
	var cfg Config
	flag.StringVar(&cfg.Host, "host", defaultHost, "flag for configure host")
	flag.IntVar(&cfg.Port, "port", defaultPort, "flag for configure port")
	flag.BoolVar(&cfg.Debug, "debug", false, "enable debug logger level")
	flag.StringVar(&cfg.DBConnStr, "db", defaultDB, "flag for configure db connection string")

	flag.Parse()

	if cfg.Host == defaultHost {
		cfg.Host = cmp.Or(os.Getenv("NOTES_HOST"), defaultHost)
	}

	if cfg.Port == defaultPort {
		port := cmp.Or(os.Getenv("NOTES_PORT"), strconv.Itoa(defaultPort))
		portInt, err := strconv.Atoi(port)
		if err != nil {
			return nil, err
		}
		cfg.Port = portInt
	}

	if cfg.DBConnStr == defaultDB {
		cfg.DBConnStr = cmp.Or(os.Getenv("NOTES_DB"), defaultDB)
	}

	return &cfg, nil
}
