FROM golang:1.15.7-alpine3.12 as builder

COPY go.mod go.sum /go/src/gitlab.com/idoko/letterpress/
WORKDIR /go/src/gitlab.com/idoko/letterpress
RUN go mod download
COPY . /go/src/gitlab.com/idoko/letterpress
RUN go build -o build/letterpress gitlab.com/idoko/letterpress/cmd/api

FROM alpine

RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/gitlab.com/idoko/letterpress/build/letterpress /usr/bin/letterpress

EXPOSE 8080 8080

ENTRYPOINT ["/usr/bin/letterpress"]