package category

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/pagination"
)

// Repository implements the category.Gateway interface using MySQL.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new category repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a new category into the database.
func (r *Repository) Create(ctx context.Context, cat *category.Category) (*category.Category, error) {
	query := `
		INSERT INTO categories (id, name, description, active, created_at, updated_at, deleted_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	var deletedAt *time.Time
	if !cat.IsActive() {
		deletedAt = new(cat.GetUpdatedAt())
	}

	_, err := r.db.ExecContext(ctx, query,
		cat.GetId().GetValue(),
		cat.GetName(),
		cat.GetDescription(),
		cat.IsActive(),
		cat.GetCreatedAt(),
		cat.GetUpdatedAt(),
		deletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}

	return cat, nil
}

// DeleteById marks a category as deleted (soft delete).
func (r *Repository) DeleteById(ctx context.Context, id category.Id) error {
	query := `
		UPDATE categories
		SET active = false, updated_at = ?, deleted_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query, now, now, id.GetValue())
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete category: get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("category not found")
	}

	return nil
}

// FindById retrieves a category by its ID.
func (r *Repository) FindById(ctx context.Context, id category.Id) (*category.Category, error) {
	query := `
		SELECT id, name, description, active, created_at, updated_at, deleted_at
		FROM categories
		WHERE id = ?
	`

	var name string
	var description sql.NullString
	var active bool
	var createdAt, updatedAt time.Time
	var deletedAt *time.Time
	var dbID string

	err := r.db.QueryRowContext(ctx, query, id.GetValue()).
		Scan(&dbID, &name, &description, &active, &createdAt, &updatedAt, &deletedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("category not found")
		}
		return nil, fmt.Errorf("find category by id: %w", err)
	}

	catID, err := category.FromString(dbID)
	if err != nil {
		return nil, fmt.Errorf("invalid category id: %w", err)
	}

	return mapToCategory(catID, name, nullableStringValue(description), active, createdAt, updatedAt, deletedAt)
}

// Update modifies an existing category.
func (r *Repository) Update(ctx context.Context, cat *category.Category) error {
	query := `
		UPDATE categories
		SET name = ?, description = ?, active = ?, updated_at = ?, deleted_at = ?
		WHERE id = ?
	`

	deletedAt := cat.GetDeletedAt()

	result, err := r.db.ExecContext(ctx, query,
		cat.GetName(),
		cat.GetDescription(),
		cat.IsActive(),
		cat.GetUpdatedAt(),
		deletedAt,
		cat.GetId().GetValue(),
	)
	if err != nil {
		return fmt.Errorf("update category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update category: get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("category not found")
	}

	return nil
}

// FindAll retrieves a paginated list of categories with optional search.
func (r *Repository) FindAll(ctx context.Context, query pagination.SearchQuery) (pagination.Pagination[category.Category], error) {
	// Build the WHERE clause
	whereClauses := []string{"deleted_at IS NULL"}
	var args []any

	if query.Terms != "" {
		searchTerm := fmt.Sprintf("%%%s%%", query.Terms)
		args = append(args, searchTerm, searchTerm)
		whereClauses = append(whereClauses, "(name LIKE ? OR description LIKE ?)")
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Determine sort column and direction
	sortColumn := sanitizeSortColumn(query.Sort)
	direction := sanitizeDirection(query.Direction)

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM categories %s", whereSQL)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return pagination.Pagination[category.Category]{}, fmt.Errorf("count categories: %w", err)
	}

	// Fetch page
	offset := (query.Page - 1) * query.PerPage
	dataQuery := fmt.Sprintf(
		"SELECT id, name, description, active, created_at, updated_at, deleted_at FROM categories %s ORDER BY %s %s LIMIT ? OFFSET ?",
		whereSQL, sortColumn, direction,
	)
	args = append(args, query.PerPage, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return pagination.Pagination[category.Category]{}, fmt.Errorf("find all categories: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("failed to close rows: %v\n", err)
		}
	}(rows)

	var items []category.Category
	for rows.Next() {
		var name string
		var description sql.NullString
		var active bool
		var createdAt, updatedAt time.Time
		var deletedAt *time.Time
		var dbID string

		if err := rows.Scan(&dbID, &name, &description, &active, &createdAt, &updatedAt, &deletedAt); err != nil {
			return pagination.Pagination[category.Category]{}, fmt.Errorf("scan category row: %w", err)
		}

		catID, err := category.FromString(dbID)
		if err != nil {
			return pagination.Pagination[category.Category]{}, fmt.Errorf("invalid category id: %w", err)
		}

		cat, err := mapToCategory(catID, name, nullableStringValue(description), active, createdAt, updatedAt, deletedAt)
		if err != nil {
			return pagination.Pagination[category.Category]{}, err
		}
		items = append(items, *cat)
	}

	if err := rows.Err(); err != nil {
		return pagination.Pagination[category.Category]{}, fmt.Errorf("iterate categories: %w", err)
	}

	if items == nil {
		items = []category.Category{}
	}

	return pagination.Pagination[category.Category]{
		CurrentPage: query.Page,
		PerPage:     query.PerPage,
		Total:       total,
		Items:       items,
	}, nil
}

// ExistsByIDs checks which category IDs exist in the database.
func (r *Repository) ExistsByIDs(ctx context.Context, ids []category.Id) ([]category.Id, error) {
	if len(ids) == 0 {
		return []category.Id{}, nil
	}

	// Build placeholder list
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id.GetValue()
	}

	query := fmt.Sprintf(
		"SELECT id FROM categories WHERE id IN (%s) AND deleted_at IS NULL",
		strings.Join(placeholders, ", "),
	)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exists by ids: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("failed to close rows: %v\n", err)
		}
	}(rows)

	var found []category.Id
	for rows.Next() {
		var dbID string
		if err := rows.Scan(&dbID); err != nil {
			return nil, fmt.Errorf("scan category id: %w", err)
		}
		catID, err := category.FromString(dbID)
		if err != nil {
			return nil, fmt.Errorf("invalid category id: %w", err)
		}
		found = append(found, catID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate category ids: %w", err)
	}

	if found == nil {
		found = []category.Id{}
	}

	return found, nil
}

// mapToCategory creates a Category entity from database values.
func mapToCategory(id category.Id, name, description string, active bool, createdAt, updatedAt time.Time, deletedAt *time.Time) (*category.Category, error) {
	cat := category.FromDB(id, name, description, active, createdAt, updatedAt, deletedAt)
	return cat, nil
}

func nullableStringValue(value sql.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}

// sanitizeSortColumn ensures the sort column is valid.
func sanitizeSortColumn(sort string) string {
	allowed := map[string]bool{
		"name":        true,
		"description": true,
		"created_at":  true,
		"updated_at":  true,
	}
	if sort == "" || !allowed[sort] {
		return "name"
	}
	return sort
}

// sanitizeDirection ensures the sort direction is valid.
func sanitizeDirection(direction string) string {
	if direction == "" || (strings.ToUpper(direction) != "ASC" && strings.ToUpper(direction) != "DESC") {
		return "ASC"
	}
	return strings.ToUpper(direction)
}
