FROM golang:1.25-alpine AS build
WORKDIR /go/src/messenger_v2

RUN apk add --no-cache git

# Копируем ВСЁ сразу (включая go.mod, go.sum и все папки)
COPY . .

# Скачиваем зависимости (если есть go.mod)
RUN go mod download || true

# Собираем (CGO отключаем для Alpine)
RUN CGO_ENABLED=0 go build -o /out/messenger ./cmd/server

FROM alpine:3.19
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=build /out/messenger /app/messenger
COPY web /app/web

EXPOSE 8020

ENV DB_PASSWORD=rootpass

ENTRYPOINT ["/app/messenger"]