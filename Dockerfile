FROM golang:1.25-alpine AS build
WORKDIR /go/src/messenger_v2

RUN apk add --no-cache git

COPY . .
ENV GOPATH=/go
ENV GO111MODULE=auto
RUN go get github.com/go-sql-driver/mysql && \
    go get github.com/gorilla/sessions && \
    go get golang.org/x/crypto/bcrypt && \
    go build -o /out/messenger ./cmd/server

FROM alpine:3.19
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=build /out/messenger /app/messenger
COPY web /app/web

EXPOSE 8020

ENV DB_PASSWORD=rootpass

ENTRYPOINT ["/app/messenger"]
