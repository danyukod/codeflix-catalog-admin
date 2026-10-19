-- Migration: Create categories table
-- Description: Creates the categories table for storing video categories

CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Index for faster searches by name
CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);

-- Index for faster searches by description
CREATE INDEX IF NOT EXISTS idx_categories_description ON categories(description);

-- Index for filtering active categories
CREATE INDEX IF NOT EXISTS idx_categories_active ON categories(active);

-- Index for soft delete filtering
CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories(deleted_at);
