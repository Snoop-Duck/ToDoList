package internal

import (
	"cmp"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Host      string `json:"-"`
	Port      int    `json:"-"`
	Debug     bool   `json:"debug"`
	DBConnStr string `json:"-"`

	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	HTTPPort   string `json:"http_port"`
}

const (
	defaultHost = "0.0.0.0"
	defaultPort = 8080
	defaultDB   = "postgres://user.password@localhost:5432/notes?sslmode=disable"
)

func ReadConfig() (*Config, error) {
	var cfg Config
	var configFile string

	flag.StringVar(&configFile, "config", "", "path to config file")
	flag.StringVar(&configFile, "c", "", "path to config file (shorthand)")

	flag.Parse()

	if configFile == "" {
		configFile = os.Getenv("CONFIG")
	}

	if configFile != "" {
		if err := loadConfigFromFile(configFile, &cfg); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}

	cfg.buildDBConnStrFromComponents()

	flag.StringVar(&cfg.Host, "host", cfg.Host, "flag for configure host")
	flag.IntVar(&cfg.Port, "port", cfg.Port, "flag for configure port")
	flag.BoolVar(&cfg.Debug, "debug", cfg.Debug, "enable debug logger level")
	flag.StringVar(&cfg.DBConnStr, "db", cfg.DBConnStr, "flag for configure db connection string")

	flag.Parse()

	applyEnvironmentOverrides(&cfg)

	return &cfg, nil
}

func (c *Config) buildDBConnStrFromComponents() {
	if c.DBConnStr != "" {
		return
	}

	if c.DBUser != "" && c.DBPassword != "" && c.DBName != "" && c.DBHost != "" {
		dbPort := cmp.Or(c.DBPort, "5432")
		c.DBConnStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			c.DBUser, c.DBPassword, c.DBHost, dbPort, c.DBName)
	}

	if c.HTTPPort != "" && c.Port == 0 {
		if port, err := strconv.Atoi(c.HTTPPort); err == nil {
			c.Port = port
		}
	}

	if c.Host == "" {
		c.Host = defaultHost
	}
}

func loadConfigFromFile(filename string, cfg *Config) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return fmt.Errorf("failed to decode config JSON: %w", err)
	}

	return nil
}

func applyEnvironmentOverrides(cfg *Config) {
	if envHost := os.Getenv("NOTES_HOST"); envHost != "" {
		cfg.Host = envHost
	} else if cfg.Host == "" {
		cfg.Host = defaultHost
	}

	if envPort := os.Getenv("NOTES_PORT"); envPort != "" {
		if port, err := strconv.Atoi(envPort); err == nil {
			cfg.Port = port
		}
	} else if cfg.Port == 0 {
		cfg.Port = defaultPort
	}

	if envDebug := os.Getenv("NOTES_DEBUG"); envDebug != "" {
		if debug, err := strconv.ParseBool(envDebug); err == nil {
			cfg.Debug = debug
		}
	}

	if envDB := os.Getenv("NOTES_DB"); envDB != "" {
		cfg.DBConnStr = envDB
	} else if cfg.DBConnStr == "" {
		cfg.DBConnStr = defaultDB
	}

	if envDBUser := os.Getenv("DB_USER"); envDBUser != "" {
		cfg.DBUser = envDBUser
		cfg.buildDBConnStrFromComponents()
	}
	if envDBPass := os.Getenv("DB_PASSWORD"); envDBPass != "" {
		cfg.DBPassword = envDBPass
		cfg.buildDBConnStrFromComponents()
	}
	if envDBName := os.Getenv("DB_NAME"); envDBName != "" {
		cfg.DBName = envDBName
		cfg.buildDBConnStrFromComponents()
	}
	if envDBHost := os.Getenv("DB_HOST"); envDBHost != "" {
		cfg.DBHost = envDBHost
		cfg.buildDBConnStrFromComponents()
	}
	if envDBPort := os.Getenv("DB_PORT"); envDBPort != "" {
		cfg.DBPort = envDBPort
		cfg.buildDBConnStrFromComponents()
	}
}
