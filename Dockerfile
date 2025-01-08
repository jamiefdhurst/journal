FROM golang:1.22-bookworm

WORKDIR /go/src/github.com/jamiefdhurst/journal
COPY . .

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends --assume-yes build-essential libsqlite3-dev; \
    go mod download; \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -v -o journal .; \
    mv journal /go/bin/journal

FROM debian:bookworm
LABEL org.opencontainers.image.source=https://github.com/jamiefdhurst/journal

WORKDIR /go/src/github.com/jamiefdhurst/journal
COPY --from=0 /go/bin/journal /usr/local/bin/
COPY --from=0 /go/src/github.com/jamiefdhurst/journal/web web

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends --assume-yes libsqlite3-0

ENV GOPATH "/go"
ENV J_ARTICLES_PER_PAGE ""
ENV J_CREATE ""
ENV J_DB_PATH ""
ENV J_DESCRIPTION ""
ENV J_EDIT ""
ENV J_GA_CODE ""
ENV J_PORT ""
ENV J_THEME ""
ENV J_TITLE ""

VOLUME /go/data
EXPOSE 3000

CMD ["journal"]
