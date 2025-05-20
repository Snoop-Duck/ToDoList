package internal

import "flag"

type Config struct {
	Host  string
	Port  int
	Debug bool
}

func ReadConfig() *Config {
	var cfg Config
	flag.StringVar(&cfg.Host, "host", "0.0.0.0", "flag for configure host")
	flag.IntVar(&cfg.Port, "port", 8080, "flag for configure port")
	flag.BoolVar(&cfg.Debug, "debug", false, "enable debug logger level")
	flag.Parse()

	return &cfg
}
