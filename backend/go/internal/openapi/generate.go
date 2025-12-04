// Code generation entrypoint for OpenAPI -> Go
package openapi

//go:generate oapi-codegen -generate types,server -package openapi -o openapi.gen.go ../../openapi.yaml

// The generated file `openapi.gen.go` will contain typed models and server interfaces
// generated from `backend/go/openapi.yaml` using `oapi-codegen`.

// To run generation locally (install oapi-codegen first):
//   go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
//   cd backend/go/internal/openapi
//   go generate
