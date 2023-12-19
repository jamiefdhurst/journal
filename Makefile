.PHONY: test

test:
	@2>&1 go test -coverprofile=cover.out -coverpkg=./internal/...,./pkg/... -v ./... | go2xunit > tests.xml
	@gocov convert cover.out | gocov-xml > coverage.xml