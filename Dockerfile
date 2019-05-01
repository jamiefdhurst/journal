FROM golang:latest

RUN apt-get update -y && apt-get install -y sqlite

WORKDIR /go/src/github.com/jamiefdhurst/journal
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENV GIPHY_API_KEY ""

VOLUME /go/data
EXPOSE 3000

CMD ["journal"]