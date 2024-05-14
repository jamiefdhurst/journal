.PHONY: build test

build:
	@CC=x86_64-unknown-linux-gnu-gcc CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -v -o bootstrap .
	@zip -r lambda.zip bootstrap web -x web/app/\*

test:
	@2>&1 go test -coverprofile=cover.out -coverpkg=./internal/...,./pkg/... -v ./... | go2xunit > tests.xml
	@gocov convert cover.out | gocov-xml > coverage.xml

clean:
	@rm -f bootstrap
	@rm -f lambda.zip