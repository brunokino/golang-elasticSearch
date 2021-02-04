FROM golang:1.15.7-buster

RUN go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate

COPY go.mod go.sum /go/src/gitlab.com/idoko/letterpress/
WORKDIR /go/src/gitlab.com/idoko/letterpress
RUN go mod download
COPY . /go/src/gitlab.com/idoko/letterpress
RUN go build -o /usr/bin/letterpress gitlab.com/idoko/letterpress/cmd/api

EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/letterpress"]