package category

func categoryURL(id string) string { return "/categories/" + id }

func categoryLinks(id string) map[string]string {
	return map[string]string{"self": categoryURL(id), "collection": "/categories"}
}
