-- Migration: Create categories table
-- Description: Creates the categories table for storing video categories

CREATE TABLE IF NOT EXISTS categories (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

-- Index for faster searches by name
CREATE INDEX idx_categories_name ON categories(name);

-- Index for faster searches by description
CREATE INDEX idx_categories_description ON categories(description(255));

-- Index for filtering active categories
CREATE INDEX idx_categories_active ON categories(active);

-- Index for soft delete filtering
CREATE INDEX idx_categories_deleted_at ON categories(deleted_at);
