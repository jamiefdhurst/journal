FROM golang:latest

RUN apt-get update -y && apt-get install -y sqlite3

WORKDIR /go/src/github.com/jamiefdhurst/journal
COPY . .

RUN go get -v ./...
RUN go install -v ./...

ENV J_ARTICLES_PER_PAGE ""
ENV J_DB_PATH ""
ENV J_GIPHY_API_KEY ""
ENV J_PORT ""
ENV J_TITLE ""

VOLUME /go/data
EXPOSE 3000

CMD ["journal"]
