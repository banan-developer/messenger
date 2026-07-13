# ==========================================
# 1. СТАДИЯ РАЗРАБОТКИ (Hot-Reload с Air)
# ==========================================
FROM golang:1.25-alpine AS dev
WORKDIR /app

RUN apk add --no-cache git && \
    go install github.com/air-verse/air@latest

COPY . .

# Инициализируем модуль (если его нет) и скачиваем зависимости
RUN (test -f go.mod || go mod init messenger_v2) && go mod tidy

EXPOSE 8020

# Запускаем Air. Перед стартом делаем tidy на случай, 
# если через Volume добавились новые импорты в код
CMD ["sh", "-c", "(test -f go.mod || go mod init messenger_v2) && go mod tidy && air -c .air.toml"]

# ==========================================
# 2. СТАДИЯ СБОРКИ (Для продакшена)
# ==========================================
FROM golang:1.25-alpine AS build
WORKDIR /go/src/messenger_v2

RUN apk add --no-cache git

COPY . .

# Инициализируем модуль внутри контейнера, затем скачиваем зависимости и собираем проект
RUN (test -f go.mod || go mod init messenger_v2) && \
    go mod tidy && \
    go build -o /out/messenger ./cmd/server

# ==========================================
# 3. ФИНАЛЬНЫЙ ОБРАЗ
# ==========================================
FROM alpine:3.19 AS prod
WORKDIR /go/src/messenger_v2

RUN apk add --no-cache ca-certificates

COPY --from=build /out/messenger /app/messenger
COPY web /app/web

EXPOSE 8020

ENV DB_PASSWORD=rootpass

ENTRYPOINT ["/app/messenger"]