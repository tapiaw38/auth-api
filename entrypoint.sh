#!/bin/sh

# Wait for PostgreSQL to be ready
until pg_isready --host=database --port=5432 --username=postgres --dbname=auth
do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

# Wait for RabbitMQ to be ready
until nc -z rabbitmq 5672
do
  echo "Waiting for RabbitMQ to be ready..."
  sleep 2
done

# Both PostgreSQL and RabbitMQ are now ready, so run the Go application
echo "========== Starting Go application =========="
exec go run cmd/main.go --host 0.0.0.0