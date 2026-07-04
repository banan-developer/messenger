FROM golang:1.25-alpine AS build
WORKDIR /src

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy && go build -o /out/messenger ./cmd/server

FROM alpine:3.19
WORKDIR /app

RUN apk add --no-cache ca-certificates netcat-openbsd socat

COPY --from=build /out/messenger /app/messenger
COPY web /app/web
COPY docker/entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh

EXPOSE 8020

ENV DB_PASSWORD=rootpass

ENTRYPOINT ["/app/entrypoint.sh"]
