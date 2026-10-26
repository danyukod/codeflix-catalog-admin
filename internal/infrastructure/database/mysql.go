package database

import (
	"cmp"
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Config holds the MySQL connection configuration.
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

// Connection wraps a sql.DB and provides database operations.
type Connection struct {
	db *sql.DB
}

// NewConnection creates a new MySQL connection pool and verifies it.
func NewConnection(ctx context.Context, cfg Config) (*Connection, error) {
	tlsConfig := cfg.SSLMode
	if tlsConfig == "disable" || tlsConfig == "" {
		tlsConfig = "false"
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?tls=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		tlsConfig,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql connection: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)

	// Verify the connection is working
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	return &Connection{db: db}, nil
}

// Pool returns the underlying sql.DB for query execution.
func (c *Connection) Pool() *sql.DB {
	return c.db
}

// Close closes the connection pool.
func (c *Connection) Close() {
	if c.db != nil {
		err := c.db.Close()
		if err != nil {
			return
		}
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
	slices.SortFunc(migrations, func(a, b migration) int {
		return cmp.Compare(a.name, b.name)
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

			if _, err := c.db.ExecContext(ctx, stmt); err != nil {
				// Ignore duplicate key errors (e.g., rerunning CREATE INDEX on existing indexes)
				if mysqlErr, ok := errors.AsType[*mysql.MySQLError](err); ok && mysqlErr.Number == 1061 {
					continue
				}
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

	for line := range strings.SplitSeq(sql, "\n") {
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
