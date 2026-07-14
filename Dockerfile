FROM golang:1.25-alpine AS build
WORKDIR /go/src/messenger_v2

RUN apk add --no-cache git

# Копируем только go.mod (без go.sum)
COPY go.mod ./
RUN go mod download

# Копируем остальной код
COPY . .

# Собираем
RUN go build -o /out/messenger ./cmd/server

FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=build /out/messenger /app/messenger
COPY web /app/web

EXPOSE 8020
ENV DB_PASSWORD=rootpass
ENTRYPOINT ["/app/messenger"]