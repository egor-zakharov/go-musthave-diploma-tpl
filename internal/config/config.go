package config

import (
	"flag"
	"os"
)

type Config struct {
	FlagRunAddr  string
	FlagDB       string
	FlagAccAddr  string
	FlagLogLevel string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) ParseFlag() {
	flag.StringVar(&c.FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&c.FlagDB, "d", "postgres://postgres:admin@localhost:5432/diploma?sslmode=disable", "database dsn")
	flag.StringVar(&c.FlagAccAddr, "r", "http://localhost:8081", "accrual system address")
	flag.StringVar(&c.FlagLogLevel, "l", "debug", "log level")

	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
		c.FlagRunAddr = envRunAddr
	}

	if envDB := os.Getenv("DATABASE_URI"); envDB != "" {
		c.FlagDB = envDB
	}

	if envAccAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccAddr != "" {
		c.FlagAccAddr = envAccAddr
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		c.FlagLogLevel = envLogLevel
	}

}
