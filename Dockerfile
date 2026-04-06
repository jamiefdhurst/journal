FROM golang:1.24-bookworm

WORKDIR /go/src/github.com/jamiefdhurst/journal
COPY . .

RUN go mod download; \
    CGO_ENABLED=0 go build -ldflags="-w -s" -o journal ./cmd/journal; \
    mv journal /go/bin/journal

FROM debian:bookworm
LABEL org.opencontainers.image.source=https://github.com/jamiefdhurst/journal

WORKDIR /go/src/github.com/jamiefdhurst/journal
COPY --from=0 /go/bin/journal /usr/local/bin/
COPY --from=0 /go/src/github.com/jamiefdhurst/journal/api api
COPY --from=0 /go/src/github.com/jamiefdhurst/journal/web web

ENV GOPATH "/go"
ENV J_CREATE ""
ENV J_DB_PATH ""
ENV J_DESCRIPTION ""
ENV J_EDIT ""
ENV J_GA_CODE ""
ENV J_PORT ""
ENV J_POSTS_PER_PAGE ""
ENV J_THEME ""
ENV J_TITLE ""

VOLUME /go/data
EXPOSE 3000

CMD ["journal"]
