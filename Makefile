.PHONY: build test coverage clean

build:
	@go build -v -o journal ./cmd/journal

test:
	@2>&1 go test -coverprofile=cover.out -coverpkg=./internal/...,./pkg/... -v ./... | go2xunit > tests.xml
	@gocov convert cover.out | gocov-xml > coverage.xml

coverage:
	@go test -coverprofile=cover.out -coverpkg=./internal/...,./pkg/... ./...
	@go tool cover -func=cover.out

clean:
	@rm -f journal