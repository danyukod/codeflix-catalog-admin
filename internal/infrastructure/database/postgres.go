package database

import (
	"context"
	"embed"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Config holds the PostgreSQL connection configuration.
type Config struct {
	Host                  string
	Port                  int
	Name                  string
	User                  string
	Password              string
	SSLMode               string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
}

// Connection wraps a pgxpool.Pool and provides database operations.
type Connection struct {
	pool *pgxpool.Pool
}

// NewConnection creates a new PostgreSQL connection pool and verifies it.
func NewConnection(ctx context.Context, cfg Config) (*Connection, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse postgres dsn: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxOpenConnections)
	poolConfig.MinConns = int32(cfg.MaxIdleConnections)
	poolConfig.MaxConnLifetime = cfg.ConnectionMaxLifetime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	// Verify the connection is working
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &Connection{pool: pool}, nil
}

// Pool returns the underlying pgxpool.Pool for query execution.
func (c *Connection) Pool() *pgxpool.Pool {
	return c.pool
}

// Close closes the connection pool.
func (c *Connection) Close() {
	if c.pool != nil {
		c.pool.Close()
	}
}

// RunMigrations executes all pending SQL migrations.
func (c *Connection) RunMigrations(ctx context.Context) error {
	// Read all migration files
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	// Sort migrations by name to ensure correct order
	type migration struct {
		name string
		sql  string
	}

	var migrations []migration
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		content, err := migrationsFS.ReadFile("migrations/" + entry.Name())
		if err != nil {
			return fmt.Errorf("read migration %s: %w", entry.Name(), err)
		}

		migrations = append(migrations, migration{
			name: entry.Name(),
			sql:  string(content),
		})
	}

	// Sort by name
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].name < migrations[j].name
	})

	// Execute each migration
	for _, m := range migrations {
		// Split by statements and execute each one
		statements := splitSQLStatements(m.sql)
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			if _, err := c.pool.Exec(ctx, stmt); err != nil {
				return fmt.Errorf("execute migration %s: %w", m.name, err)
			}
		}
	}

	return nil
}

// splitSQLStatements splits SQL content by semicolons.
func splitSQLStatements(sql string) []string {
	var statements []string
	var current strings.Builder

	for _, line := range strings.Split(sql, "\n") {
		trimmed := strings.TrimSpace(line)
		// Skip comments
		if strings.HasPrefix(trimmed, "--") {
			continue
		}
		current.WriteString(line)
		current.WriteString("\n")
		if strings.HasSuffix(trimmed, ";") {
			statements = append(statements, current.String())
			current.Reset()
		}
	}

	// Add remaining content
	if current.Len() > 0 {
		statements = append(statements, current.String())
	}

	return statements
}
