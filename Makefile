OPENAPI_VALIDATE_VERSION := 0.123.0

run:
	go run main.go

validate-api:
	go run github.com/getkin/kin-openapi/cmd/validate@v$(OPENAPI_VALIDATE_VERSION) -- api-spec.yaml

tidy:
	go mod tidy
	go mod vendor
