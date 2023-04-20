.PHONY: test

test:
	@2>&1 go test -coverprofile=cover.out -coverpkg=./internal/...,./pkg/... -v ./... | go2xunit
	@gocover-cobertura < cover.out > coverage.xml