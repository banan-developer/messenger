#!/bin/sh
set -eu

echo "Waiting for database..."
until nc -z db 3306; do
  sleep 1
done

echo "Starting local MySQL proxy..."
socat TCP-LISTEN:3306,fork,reuseaddr TCP:db:3306 &
PROXY_PID=$!

trap 'kill $PROXY_PID >/dev/null 2>&1 || true' EXIT

echo "Starting messenger server..."
exec /app/messenger
