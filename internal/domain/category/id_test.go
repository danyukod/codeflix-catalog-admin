package category

import "testing"

func TestFromStringRequiresUUID(t *testing.T) {
	for _, id := range []string{"", "not-a-uuid"} {
		t.Run(id, func(t *testing.T) {
			if _, err := FromString(id); err == nil {
				t.Fatalf("FromString(%q) expected an error", id)
			}
		})
	}

	created := NewCategory("Movies", "", true)
	if _, err := FromString(created.GetId().GetValue()); err != nil {
		t.Fatalf("FromString(valid UUID): %v", err)
	}
}
