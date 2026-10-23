package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_splitSQLStatements(t *testing.T) {
	t.Run("given single statement when splitting should return one statement", func(t *testing.T) {
		sql := "CREATE TABLE users (id INT);"
		statements := splitSQLStatements(sql)

		assert.Len(t, statements, 1)
		assert.Contains(t, statements[0], "CREATE TABLE users")
	})

	t.Run("given multiple statements when splitting should return all statements", func(t *testing.T) {
		sql := "CREATE TABLE users (id INT);\nCREATE TABLE posts (id INT);"
		statements := splitSQLStatements(sql)

		assert.Len(t, statements, 2)
		assert.Contains(t, statements[0], "CREATE TABLE users")
		assert.Contains(t, statements[1], "CREATE TABLE posts")
	})

	t.Run("given SQL with comments when splitting should skip comments", func(t *testing.T) {
		sql := "-- Create users table\nCREATE TABLE users (id INT);"
		statements := splitSQLStatements(sql)

		assert.Len(t, statements, 1)
		assert.NotContains(t, statements[0], "-- Create users table")
		assert.Contains(t, statements[0], "CREATE TABLE users")
	})

	t.Run("given SQL with multiple comments when splitting should skip all comments", func(t *testing.T) {
		sql := "-- Comment 1\n-- Comment 2\nCREATE TABLE users (id INT);\n-- Comment 3"
		statements := splitSQLStatements(sql)

		assert.Len(t, statements, 1)
		assert.Contains(t, statements[0], "CREATE TABLE users")
	})

	t.Run("given empty SQL when splitting should return single newline", func(t *testing.T) {
		sql := ""
		statements := splitSQLStatements(sql)

		// splitSQLStatements splits by "\n", so empty string produces one line "\n"
		assert.Len(t, statements, 1)
		assert.Equal(t, "\n", statements[0])
	})

	t.Run("given SQL with blank lines when splitting should handle correctly", func(t *testing.T) {
		sql := "CREATE TABLE users (id INT);\n\nCREATE TABLE posts (id INT);"
		statements := splitSQLStatements(sql)

		assert.Len(t, statements, 2)
	})

	t.Run("given SQL without semicolon when splitting should return statement", func(t *testing.T) {
		sql := "CREATE TABLE users (id INT)"
		statements := splitSQLStatements(sql)

		assert.Len(t, statements, 1)
		assert.Contains(t, statements[0], "CREATE TABLE users")
	})

	t.Run("given SQL with only comments when splitting should return empty slice", func(t *testing.T) {
		sql := "-- Comment 1\n-- Comment 2"
		statements := splitSQLStatements(sql)

		assert.Len(t, statements, 0)
	})

	t.Run("given SQL with multiline statement when splitting should preserve newlines", func(t *testing.T) {
		sql := "CREATE TABLE users (\n    id INT,\n    name VARCHAR(255)\n);"
		statements := splitSQLStatements(sql)

		assert.Len(t, statements, 1)
		assert.Contains(t, statements[0], "CREATE TABLE users")
		assert.Contains(t, statements[0], "id INT")
		assert.Contains(t, statements[0], "name VARCHAR(255)")
	})

	t.Run("given SQL with mixed comments and statements when splitting should handle correctly", func(t *testing.T) {
		sql := "-- Migration 1\nCREATE TABLE users (id INT);\n-- Migration 2\nCREATE TABLE posts (id INT);"
		statements := splitSQLStatements(sql)

		assert.Len(t, statements, 2)
		assert.Contains(t, statements[0], "CREATE TABLE users")
		assert.Contains(t, statements[1], "CREATE TABLE posts")
	})
}

func TestConfig(t *testing.T) {
	t.Run("given valid config when created should have correct values", func(t *testing.T) {
		cfg := Config{
			Host:               "localhost",
			Port:               3306,
			Name:               "shop",
			User:               "shop",
			Password:           "shop",
			SSLMode:            "disable",
			MaxOpenConnections: 25,
			MaxIdleConnections: 5,
		}

		assert.Equal(t, "localhost", cfg.Host)
		assert.Equal(t, 3306, cfg.Port)
		assert.Equal(t, "shop", cfg.Name)
		assert.Equal(t, "shop", cfg.User)
		assert.Equal(t, "shop", cfg.Password)
		assert.Equal(t, "disable", cfg.SSLMode)
		assert.Equal(t, 25, cfg.MaxOpenConnections)
		assert.Equal(t, 5, cfg.MaxIdleConnections)
	})
}

func TestNewConnection(t *testing.T) {
	t.Run("given invalid host when creating connection should return error", func(t *testing.T) {
		ctx := context.Background()
		cfg := Config{
			Host:               "nonexistent-host-12345.invalid",
			Port:               3306,
			Name:               "shop",
			User:               "shop",
			Password:           "shop",
			SSLMode:            "disable",
			MaxOpenConnections: 25,
			MaxIdleConnections: 5,
		}

		_, err := NewConnection(ctx, cfg)

		assert.Error(t, err)
	})

	t.Run("given invalid port when creating connection should return error", func(t *testing.T) {
		ctx := context.Background()
		cfg := Config{
			Host:               "localhost",
			Port:               99999,
			Name:               "shop",
			User:               "shop",
			Password:           "shop",
			SSLMode:            "disable",
			MaxOpenConnections: 25,
			MaxIdleConnections: 5,
		}

		_, err := NewConnection(ctx, cfg)

		assert.Error(t, err)
	})
}

func TestConnection_Pool(t *testing.T) {
	t.Run("given closed connection when calling pool should return nil", func(t *testing.T) {
		conn := &Connection{db: nil}

		pool := conn.Pool()

		assert.Nil(t, pool)
	})
}

func TestConnection_Close(t *testing.T) {
	t.Run("given nil database when closing should not panic", func(t *testing.T) {
		conn := &Connection{db: nil}

		assert.NotPanics(t, func() {
			conn.Close()
		})
	})
}

func TestRunMigrations(t *testing.T) {
	t.Run("given nil database when running migrations should panic", func(t *testing.T) {
		ctx := context.Background()
		conn := &Connection{db: nil}

		assert.Panics(t, func() {
			conn.RunMigrations(ctx)
		})
	})
}
