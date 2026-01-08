package category

import "testing"

// ... existing code ...
func TestNewCategory(t *testing.T) {
	category := Category{
		Id: "1",
	}

	if category.Id != "1" {
		t.Errorf("expected category Id to be '1', got '%s'", category.Id)
	}
}
