package docs

import _ "embed"

//go:embed openapi.yaml
var OpenAPIYAML []byte

//go:embed swagger-ui.html
var SwaggerUIHTML []byte
