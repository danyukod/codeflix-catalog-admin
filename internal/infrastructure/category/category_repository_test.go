package category

import (
	"database/sql"
	"testing"
	"time"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite"
)

// setupTestDB creates an in-memory SQLite database with the categories table.
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE categories (
			id VARCHAR(36) PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			active BOOLEAN NOT NULL DEFAULT 1,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP
		)
	`)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := db.Close()
		if err != nil {
			return
		}
	})

	return db
}

// createTestCategory creates a valid active category for testing.
func createTestCategory(t *testing.T, name, description string) *category.Category {
	t.Helper()
	cat := category.NewCategory(name, description, true)
	require.NotNil(t, cat)
	return cat
}

func TestRepository_Create(t *testing.T) {
	t.Run("given valid category when creating should insert and return category", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		cat := createTestCategory(t, "Electronics", "Electronic devices and gadgets")
		result, err := repo.Create(t.Context(), cat)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, cat.GetId().GetValue(), result.GetId().GetValue())
		assert.Equal(t, "Electronics", result.GetName())

		// Verify in database
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM categories WHERE id = ?", cat.GetId().GetValue()).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("given inactive category when creating should set deleted_at", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		cat := category.NewCategory("Books", "Physical books", false)
		_, err := repo.Create(t.Context(), cat)
		require.NoError(t, err)

		var deletedAt *time.Time
		err = db.QueryRow("SELECT deleted_at FROM categories WHERE id = ?", cat.GetId().GetValue()).Scan(&deletedAt)
		require.NoError(t, err)
		assert.NotNil(t, deletedAt)
	})

	t.Run("given duplicate ID when creating should return error", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		cat := createTestCategory(t, "Electronics", "Electronic devices")
		_, err := repo.Create(t.Context(), cat)
		require.NoError(t, err)

		// Try to insert the same category again (duplicate primary key)
		_, err = repo.Create(t.Context(), cat)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create category")
	})
}

func TestRepository_FindById(t *testing.T) {
	t.Run("given existing ID when finding should return category", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		cat := createTestCategory(t, "Movies", "Film and movie categories")
		_, err := repo.Create(t.Context(), cat)
		require.NoError(t, err)

		result, err := repo.FindById(t.Context(), cat.GetId().(category.Id))

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, cat.GetId().GetValue(), result.GetId().GetValue())
		assert.Equal(t, "Movies", result.GetName())
		assert.Equal(t, "Film and movie categories", result.GetDescription())
		assert.True(t, result.IsActive())
	})

	t.Run("given non-existent ID when finding should return error", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		nonExistentID := category.Unique()
		_, err := repo.FindById(t.Context(), nonExistentID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("given category with null description when finding should return empty description", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)
		cat := createTestCategory(t, "Without Description", "")

		_, err := db.Exec(`
			INSERT INTO categories (id, name, description, active, created_at, updated_at)
			VALUES (?, ?, NULL, ?, ?, ?)
		`, cat.GetId().GetValue(), cat.GetName(), cat.IsActive(), cat.GetCreatedAt(), cat.GetUpdatedAt())
		require.NoError(t, err)

		result, err := repo.FindById(t.Context(), cat.GetId().(category.Id))

		require.NoError(t, err)
		assert.Empty(t, result.GetDescription())
	})

	t.Run("given soft-deleted category when finding should still return it", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		cat := createTestCategory(t, "Deleted", "Soft deleted category")
		_, err := repo.Create(t.Context(), cat)
		require.NoError(t, err)

		// Soft delete the category
		now := time.Now().UTC()
		_, err = db.Exec("UPDATE categories SET active = 0, deleted_at = ? WHERE id = ?", now, cat.GetId().GetValue())
		require.NoError(t, err)

		result, err := repo.FindById(t.Context(), cat.GetId().(category.Id))

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Deleted", result.GetName())
	})
}

func TestRepository_Update(t *testing.T) {
	t.Run("given existing category when updating should modify fields", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		cat := createTestCategory(t, "Original Name", "Original description")
		_, err := repo.Create(t.Context(), cat)
		require.NoError(t, err)

		// Update the category
		cat.Update("Updated Name", "Updated description", true)
		err = repo.Update(t.Context(), cat)

		require.NoError(t, err)

		// Verify in database
		result, err := repo.FindById(t.Context(), cat.GetId().(category.Id))
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", result.GetName())
		assert.Equal(t, "Updated description", result.GetDescription())
	})

	t.Run("given non-existent category when updating should return error", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		cat := createTestCategory(t, "NonExistent", "Does not exist")
		err := repo.Update(t.Context(), cat)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("given deactivated category when updating should set deleted_at", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		cat := createTestCategory(t, "ToDeactivate", "Will be deactivated")
		_, err := repo.Create(t.Context(), cat)
		require.NoError(t, err)

		cat.Deactivate()
		err = repo.Update(t.Context(), cat)
		require.NoError(t, err)

		result, err := repo.FindById(t.Context(), cat.GetId().(category.Id))
		require.NoError(t, err)
		assert.False(t, result.IsActive())
		assert.NotNil(t, result.GetDeletedAt())
	})
}

func TestRepository_DeleteById(t *testing.T) {
	t.Run("given existing category when deleting should soft delete", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		cat := createTestCategory(t, "ToDelete", "Will be deleted")
		_, err := repo.Create(t.Context(), cat)
		require.NoError(t, err)

		err = repo.DeleteById(t.Context(), cat.GetId().(category.Id))
		require.NoError(t, err)

		// Verify soft delete in database
		var active bool
		var deletedAt *time.Time
		err = db.QueryRow("SELECT active, deleted_at FROM categories WHERE id = ?", cat.GetId().GetValue()).Scan(&active, &deletedAt)
		require.NoError(t, err)
		assert.False(t, active)
		assert.NotNil(t, deletedAt)
	})

	t.Run("given non-existent ID when deleting should return error", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		nonExistentID := category.Unique()
		err := repo.DeleteById(t.Context(), nonExistentID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("given already deleted category when deleting should return error", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		cat := createTestCategory(t, "AlreadyDeleted", "Already soft deleted")
		_, err := repo.Create(t.Context(), cat)
		require.NoError(t, err)

		// First delete
		err = repo.DeleteById(t.Context(), cat.GetId().(category.Id))
		require.NoError(t, err)

		// Second delete should fail
		err = repo.DeleteById(t.Context(), cat.GetId().(category.Id))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})
}

func TestRepository_FindAll(t *testing.T) {
	t.Run("given empty database when finding all should return empty list", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		query := pagination.NewSearchQuery(1, 10, "", "name", "ASC")
		result, err := repo.FindAll(t.Context(), query)

		require.NoError(t, err)
		assert.Equal(t, int64(0), result.Total)
		assert.Empty(t, result.Items)
		assert.Equal(t, 1, result.CurrentPage)
		assert.Equal(t, 10, result.PerPage)
	})

	t.Run("given category with null description when finding all should return empty description", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)
		cat := createTestCategory(t, "Without Description", "")

		_, err := db.Exec(`
			INSERT INTO categories (id, name, description, active, created_at, updated_at)
			VALUES (?, ?, NULL, ?, ?, ?)
		`, cat.GetId().GetValue(), cat.GetName(), cat.IsActive(), cat.GetCreatedAt(), cat.GetUpdatedAt())
		require.NoError(t, err)

		result, err := repo.FindAll(t.Context(), pagination.NewSearchQuery(1, 10, "", "name", "ASC"))

		require.NoError(t, err)
		require.Len(t, result.Items, 1)
		assert.Empty(t, result.Items[0].GetDescription())
	})

	t.Run("given multiple categories when finding all should return paginated results", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		// Create 3 categories
		for i := 1; i <= 3; i++ {
			cat := createTestCategory(t, "Category "+string(rune('A'+i-1)), "Description "+string(rune('A'+i-1)))
			_, err := repo.Create(t.Context(), cat)
			require.NoError(t, err)
		}

		query := pagination.NewSearchQuery(1, 10, "", "name", "ASC")
		result, err := repo.FindAll(t.Context(), query)

		require.NoError(t, err)
		assert.Equal(t, int64(3), result.Total)
		assert.Len(t, result.Items, 3)
	})

	t.Run("given paginated query when finding all should return correct page", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		// Create 5 categories
		for i := 1; i <= 5; i++ {
			cat := createTestCategory(t, "Cat"+string(rune('0'+i)), "Desc"+string(rune('0'+i)))
			_, err := repo.Create(t.Context(), cat)
			require.NoError(t, err)
		}

		// Page 1, 2 per page
		query := pagination.NewSearchQuery(1, 2, "", "name", "ASC")
		result, err := repo.FindAll(t.Context(), query)

		require.NoError(t, err)
		assert.Equal(t, int64(5), result.Total)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, 1, result.CurrentPage)

		// Page 2, 2 per page
		query = pagination.NewSearchQuery(2, 2, "", "name", "ASC")
		result, err = repo.FindAll(t.Context(), query)

		require.NoError(t, err)
		assert.Equal(t, int64(5), result.Total)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, 2, result.CurrentPage)
	})

	t.Run("given search terms when finding all should filter by name or description", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		_, err := repo.Create(t.Context(), createTestCategory(t, "Electronics", "Electronic devices"))
		require.NoError(t, err)
		_, err = repo.Create(t.Context(), createTestCategory(t, "Books", "Physical books"))
		require.NoError(t, err)
		_, err = repo.Create(t.Context(), createTestCategory(t, "Electrical Tools", "Power tools"))
		require.NoError(t, err)

		query := pagination.NewSearchQuery(1, 10, "Electronic", "", "")
		result, err := repo.FindAll(t.Context(), query)

		require.NoError(t, err)
		// SQLite LIKE is case-insensitive by default
		assert.GreaterOrEqual(t, result.Total, int64(1))
		assert.NotEmpty(t, result.Items)
	})

	t.Run("given soft-deleted categories when finding all should exclude them", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		cat1 := createTestCategory(t, "Active", "Active category")
		cat2 := createTestCategory(t, "Deleted", "Deleted category")

		_, err := repo.Create(t.Context(), cat1)
		require.NoError(t, err)
		_, err = repo.Create(t.Context(), cat2)
		require.NoError(t, err)

		// Soft delete cat2
		err = repo.DeleteById(t.Context(), cat2.GetId().(category.Id))
		require.NoError(t, err)

		query := pagination.NewSearchQuery(1, 10, "", "name", "ASC")
		result, err := repo.FindAll(t.Context(), query)

		require.NoError(t, err)
		assert.Equal(t, int64(1), result.Total)
		assert.Len(t, result.Items, 1)
		assert.Equal(t, "Active", result.Items[0].GetName())
	})

	t.Run("given sort parameter when finding all should sort correctly", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		_, err := repo.Create(t.Context(), createTestCategory(t, "Zebra", "Z description"))
		require.NoError(t, err)
		_, err = repo.Create(t.Context(), createTestCategory(t, "Apple", "A description"))
		require.NoError(t, err)
		_, err = repo.Create(t.Context(), createTestCategory(t, "Mango", "M description"))
		require.NoError(t, err)

		query := pagination.NewSearchQuery(1, 10, "", "name", "ASC")
		result, err := repo.FindAll(t.Context(), query)

		require.NoError(t, err)
		assert.Len(t, result.Items, 3)
		assert.Equal(t, "Apple", result.Items[0].GetName())
		assert.Equal(t, "Mango", result.Items[1].GetName())
		assert.Equal(t, "Zebra", result.Items[2].GetName())
	})

	t.Run("given invalid sort column when finding all should default to name", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		_, err := repo.Create(t.Context(), createTestCategory(t, "B", "B desc"))
		require.NoError(t, err)
		_, err = repo.Create(t.Context(), createTestCategory(t, "A", "A desc"))
		require.NoError(t, err)

		query := pagination.NewSearchQuery(1, 10, "", "invalid_column", "ASC")
		result, err := repo.FindAll(t.Context(), query)

		require.NoError(t, err)
		assert.Len(t, result.Items, 2)
		// Should default to name ASC
		assert.Equal(t, "A", result.Items[0].GetName())
	})
}

func TestRepository_ExistsByIDs(t *testing.T) {
	t.Run("given empty IDs when checking existence should return empty list", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		result, err := repo.ExistsByIDs(t.Context(), []category.Id{})

		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("given existing IDs when checking existence should return found IDs", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		cat1 := createTestCategory(t, "Existing1", "First existing")
		cat2 := createTestCategory(t, "Existing2", "Second existing")

		_, err := repo.Create(t.Context(), cat1)
		require.NoError(t, err)
		_, err = repo.Create(t.Context(), cat2)
		require.NoError(t, err)

		nonExistentID := category.Unique()
		ids := []category.Id{cat1.GetId().(category.Id), nonExistentID, cat2.GetId().(category.Id)}
		result, err := repo.ExistsByIDs(t.Context(), ids)

		require.NoError(t, err)
		assert.Len(t, result, 2)

		foundIDs := make(map[string]bool)
		for _, id := range result {
			foundIDs[id.GetValue()] = true
		}
		assert.True(t, foundIDs[cat1.GetId().GetValue()])
		assert.True(t, foundIDs[cat2.GetId().GetValue()])
	})

	t.Run("given non-existent IDs when checking existence should return empty list", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		ids := []category.Id{category.Unique(), category.Unique()}
		result, err := repo.ExistsByIDs(t.Context(), ids)

		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("given soft-deleted category ID when checking existence should exclude it", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		cat := createTestCategory(t, "ToBeDeleted", "Will be deleted")
		_, err := repo.Create(t.Context(), cat)
		require.NoError(t, err)

		// Soft delete
		err = repo.DeleteById(t.Context(), cat.GetId().(category.Id))
		require.NoError(t, err)

		result, err := repo.ExistsByIDs(t.Context(), []category.Id{cat.GetId().(category.Id)})

		require.NoError(t, err)
		assert.Empty(t, result)
	})
}

func Test_sanitizeSortColumn(t *testing.T) {
	t.Run("given empty string when sanitizing should return name", func(t *testing.T) {
		result := sanitizeSortColumn("")
		assert.Equal(t, "name", result)
	})

	t.Run("given valid column name when sanitizing should return the column", func(t *testing.T) {
		assert.Equal(t, "name", sanitizeSortColumn("name"))
		assert.Equal(t, "description", sanitizeSortColumn("description"))
		assert.Equal(t, "created_at", sanitizeSortColumn("created_at"))
		assert.Equal(t, "updated_at", sanitizeSortColumn("updated_at"))
	})

	t.Run("given invalid column when sanitizing should return name", func(t *testing.T) {
		assert.Equal(t, "name", sanitizeSortColumn("invalid_column"))
		assert.Equal(t, "name", sanitizeSortColumn("DROP TABLE"))
	})
}

func Test_sanitizeDirection(t *testing.T) {
	t.Run("given empty string when sanitizing should return ASC", func(t *testing.T) {
		result := sanitizeDirection("")
		assert.Equal(t, "ASC", result)
	})

	t.Run("given valid direction when sanitizing should return uppercase", func(t *testing.T) {
		assert.Equal(t, "ASC", sanitizeDirection("ASC"))
		assert.Equal(t, "ASC", sanitizeDirection("asc"))
		assert.Equal(t, "DESC", sanitizeDirection("DESC"))
		assert.Equal(t, "DESC", sanitizeDirection("desc"))
	})

	t.Run("given invalid direction when sanitizing should return ASC", func(t *testing.T) {
		assert.Equal(t, "ASC", sanitizeDirection("INVALID"))
		assert.Equal(t, "ASC", sanitizeDirection("123"))
	})
}

func TestNewRepository(t *testing.T) {
	t.Run("given valid db when creating repository should return non-nil", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		assert.NotNil(t, repo)
	})
}
