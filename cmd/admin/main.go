package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Database   DatabaseConfig   `yaml:"database"`
	Logging    LoggingConfig    `yaml:"logging"`
	Pagination PaginationConfig `yaml:"pagination"`
}

type ServerConfig struct {
	Host            string          `yaml:"host"`
	Port            int             `yaml:"port"`
	ShutdownTimeout string          `yaml:"shutdown_timeout"`
	RequestTimeout  string          `yaml:"request_timeout"`
	Compression     CompressionConfig `yaml:"compression"`
	Undertow        UndertowConfig    `yaml:"undertow"`
}

type CompressionConfig struct {
	Enabled   bool     `yaml:"enabled"`
	Level     int      `yaml:"level"`
	MimeTypes []string `yaml:"mime_types"`
	MinSize   int      `yaml:"min_size"`
}

type UndertowConfig struct {
	WorkerThreads int `yaml:"worker_threads"`
	IOThreads     int `yaml:"io_threads"`
}

type DatabaseConfig struct {
	Driver                string `yaml:"driver"`
	Host                  string `yaml:"host"`
	Port                  int    `yaml:"port"`
	Name                  string `yaml:"name"`
	User                  string `yaml:"user"`
	Password              string `yaml:"password"`
	SSLMode               string `yaml:"ssl_mode"`
	MaxOpenConnections    int    `yaml:"max_open_connections"`
	MaxIdleConnections    int    `yaml:"max_idle_connections"`
	ConnectionMaxLifetime string `yaml:"connection_max_lifetime"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type PaginationConfig struct {
	DefaultPage    int `yaml:"default_page"`
	DefaultPerPage int `yaml:"default_per_page"`
	MaxPerPage     int `yaml:"max_per_page"`
}

func (s *ServerConfig) ParsedShutdownTimeout() (time.Duration, error) {
	return time.ParseDuration(s.ShutdownTimeout)
}

func (s *ServerConfig) ParsedRequestTimeout() (time.Duration, error) {
	return time.ParseDuration(s.RequestTimeout)
}

func (d *DatabaseConfig) ParsedConnectionMaxLifetime() (time.Duration, error) {
	return time.ParseDuration(d.ConnectionMaxLifetime)
}

func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config file: %w", err)
	}

	return &cfg, nil
}

func main() {
	cfg, err := loadConfig("configs/application.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	fmt.Printf("Server listening on %s:%d\n", cfg.Server.Host, cfg.Server.Port)
}
