FROM golang:1.25-alpine AS build
WORKDIR /go/src/messenger_v2

# Копируем go.mod и go.sum для скачивания зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код
COPY . .

# Собираем бинарник
RUN go build -o /out/messenger ./cmd/server

FROM alpine:3.19
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=build /out/messenger /app/messenger
COPY web /app/web

EXPOSE 8020

ENV DB_PASSWORD=rootpass

ENTRYPOINT ["/app/messenger"]