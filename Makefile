OPENAPI_VALIDATE_VERSION := 0.123.0

run:
	go run main.go

run-race:
	go run -race main.go

validate-api:
	go run github.com/getkin/kin-openapi/cmd/validate@v$(OPENAPI_VALIDATE_VERSION) -- api-spec.yaml

tidy:
	go mod tidy
	go mod vendor

test:
	go test -timeout 5s ./...

test-race:
	go test -race -timeout 10s ./...
