package docs

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestOpenAPISchemasResolve(t *testing.T) {
	var root map[string]any
	if err := yaml.Unmarshal(OpenAPIYAML, &root); err != nil {
		t.Fatal(err)
	}
	if root["openapi"] != "3.1.2" {
		t.Fatalf("unexpected OpenAPI version: %v", root["openapi"])
	}
	components, ok := root["components"].(map[string]any)
	if !ok {
		t.Fatal("missing components")
	}
	if _, ok := components["schemas"]; !ok {
		t.Fatal("missing schemas")
	}
	if count := walkRefs(t, root, root); count == 0 {
		t.Fatal("no references found")
	}
}

func walkRefs(t *testing.T, root map[string]any, value any) int {
	t.Helper()
	count := 0
	switch node := value.(type) {
	case map[string]any:
		if ref, ok := node["$ref"].(string); ok {
			count++
			if !strings.HasPrefix(ref, "#/") {
				t.Fatalf("unsupported external reference: %s", ref)
			}
			var target any = root
			for _, segment := range strings.Split(strings.TrimPrefix(ref, "#/"), "/") {
				object, ok := target.(map[string]any)
				if !ok {
					t.Fatalf("invalid reference: %s", ref)
				}
				target, ok = object[segment]
				if !ok {
					t.Fatalf("unresolved reference: %s", ref)
				}
			}
			if target == nil {
				t.Fatalf("empty reference: %s", ref)
			}
		}
		for _, child := range node {
			count += walkRefs(t, root, child)
		}
	case []any:
		for _, child := range node {
			count += walkRefs(t, root, child)
		}
	}
	return count
}
