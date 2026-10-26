package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	createcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/create"
	deletecategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/delete"
	getcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/retrieve/get"
	listcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/retrieve/list"
	updatecategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/update"
	"github.com/danyukod/codeflix-catalog-admin/internal/infrastructure/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/infrastructure/database"
	apphttp "github.com/danyukod/codeflix-catalog-admin/internal/interfaces/http"
	categorycontroller "github.com/danyukod/codeflix-catalog-admin/internal/interfaces/http/controller/category"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Database   DatabaseConfig   `yaml:"database"`
	Logging    LoggingConfig    `yaml:"logging"`
	Pagination PaginationConfig `yaml:"pagination"`
}

type ServerConfig struct {
	Host            string            `yaml:"host"`
	Port            int               `yaml:"port"`
	ShutdownTimeout string            `yaml:"shutdown_timeout"`
	RequestTimeout  string            `yaml:"request_timeout"`
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
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("failed to close config file: %v", err)
		}
	}(f)

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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Parse connection max lifetime
	maxLifetime, err := cfg.Database.ParsedConnectionMaxLifetime()
	if err != nil {
		log.Fatalf("failed to parse connection max lifetime: %v", err)
	}

	// Create database connection
	dbConfig := database.Config{
		Host:                  cfg.Database.Host,
		Port:                  cfg.Database.Port,
		Name:                  cfg.Database.Name,
		User:                  cfg.Database.User,
		Password:              cfg.Database.Password,
		SSLMode:               cfg.Database.SSLMode,
		MaxOpenConnections:    cfg.Database.MaxOpenConnections,
		MaxIdleConnections:    cfg.Database.MaxIdleConnections,
		ConnectionMaxLifetime: maxLifetime,
	}

	dbConn, err := database.NewConnection(ctx, dbConfig)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	fmt.Println("Successfully connected to database")

	// Run migrations
	if err := dbConn.RunMigrations(ctx); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	fmt.Println("Migrations executed successfully")

	// Create category repository
	repo := category.NewRepository(dbConn.Pool())
	controller := categorycontroller.NewController(
		createcategory.NewCreateCategoryUseCase(repo),
		listcategory.NewListCategoryUseCase(repo),
		getcategory.NewGetCategoryByIdUseCase(repo),
		updatecategory.NewUpdateCategoryUsecase(repo),
		deletecategory.NewDeleteCategoryUseCase(repo),
	)
	server := &http.Server{Addr: fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port), Handler: apphttp.NewRouter(controller)}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	fmt.Printf("Server listening on %s:%d\n", cfg.Server.Host, cfg.Server.Port)

	// Wait for shutdown signal
	<-ctx.Done()
	fmt.Println("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("failed to shutdown HTTP server: %v", err)
	}
}
