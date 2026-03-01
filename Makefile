.PHONY: build test

build:
	@go build -v -o journal ./cmd/journal

test:
	@2>&1 go test -coverprofile=cover.out -coverpkg=./internal/...,./pkg/... -v ./... | go2xunit > tests.xml
	@gocov convert cover.out | gocov-xml > coverage.xml

clean:
	@rm -f journal