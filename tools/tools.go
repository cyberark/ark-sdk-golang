//go:build generate

package tools

//go:generate go run ./genservices
//go:generate go run ./genapi
//go:generate gofmt -w ./../pkg/ark_api.go
//go:generate gofmt -w ./../pkg/ark_api_services.go
