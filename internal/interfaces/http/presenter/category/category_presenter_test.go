package category

import (
	"encoding/json"
	"testing"
)

func TestCategoryResponseSerializationRoundTrip(t *testing.T) {
	t.Parallel()

	response := CategoryResponse{
		ID:          "category-id",
		Name:        "Action",
		Description: "Action movies",
		IsActive:    true,
		CreatedAt:   "2026-10-30T12:00:00Z",
		UpdatedAt:   "2026-10-30T12:30:00Z",
		Links: map[string]string{
			"self":       "/categories/category-id",
			"collection": "/categories",
		},
	}

	serialized, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("marshal category response: %v", err)
	}

	var deserialized CategoryResponse
	if err := json.Unmarshal(serialized, &deserialized); err != nil {
		t.Fatalf("unmarshal category response: %v", err)
	}

	if deserialized.ID != response.ID {
		t.Errorf("ID = %q, want %q", deserialized.ID, response.ID)
	}
	if deserialized.Name != response.Name {
		t.Errorf("Name = %q, want %q", deserialized.Name, response.Name)
	}
	if deserialized.Description != response.Description {
		t.Errorf("Description = %q, want %q", deserialized.Description, response.Description)
	}
	if deserialized.IsActive != response.IsActive {
		t.Errorf("IsActive = %t, want %t", deserialized.IsActive, response.IsActive)
	}
	if deserialized.CreatedAt != response.CreatedAt || deserialized.UpdatedAt != response.UpdatedAt {
		t.Errorf("timestamps = (%q, %q), want (%q, %q)", deserialized.CreatedAt, deserialized.UpdatedAt, response.CreatedAt, response.UpdatedAt)
	}
	if len(deserialized.Links) != len(response.Links) || deserialized.Links["self"] != response.Links["self"] || deserialized.Links["collection"] != response.Links["collection"] {
		t.Errorf("Links = %#v, want %#v", deserialized.Links, response.Links)
	}

	if string(serialized) == "" || string(serialized)[0] != '{' {
		t.Fatalf("serialized response is not a JSON object: %s", serialized)
	}
}

func TestCategoryResponseSerializationOmitsDeletedAtWhenNil(t *testing.T) {
	t.Parallel()

	serialized, err := json.Marshal(CategoryResponse{ID: "category-id"})
	if err != nil {
		t.Fatalf("marshal category response: %v", err)
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(serialized, &fields); err != nil {
		t.Fatalf("unmarshal serialized response: %v", err)
	}
	if _, ok := fields["deleted_at"]; ok {
		t.Error("deleted_at should be omitted when nil")
	}
}
