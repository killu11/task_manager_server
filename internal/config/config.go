package config

import "fmt"

type Config struct {
	DB *DatabaseConfig
}

func NewConfig() (*Config, error) {
	c := &Config{
		DB: new(DatabaseConfig),
	}
	if err := c.LoadEnv(); err != nil {
		return nil, err
	}
	return c, nil
}

type DatabaseConfig struct {
	Host     string `env:"POSTGRES_HOST" def:"localhost"`
	Port     int    `env:"POSTGRES_PORT" def:"5432"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBName   string `env:"POSTGRES_DB"`
	SSLMode  string `env:"POSTGRES_SSLMODE"`
}

func (dc *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dc.User,
		dc.Password,
		dc.Host,
		dc.Port,
		dc.DBName,
		dc.SSLMode,
	)
}
