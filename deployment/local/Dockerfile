# Stage 1: Build the Go Application
FROM golang:1.19-alpine as builder

# Add Maintainer Info
LABEL maintainer="tapiaw38 Singh <tapiaw38@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum /app/

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Install PostgreSQL client and Golang-Migrate
RUN apk update && apk add --no-cache postgresql-client curl && \
    curl -L -o /usr/local/bin/migrate https://github.com/golang-migrate/migrate/releases/download/v4.15.0/migrate.linux-amd64.tar.gz && \
    chmod +x /usr/local/bin/migrate

# Copy the rest of the application code
COPY . /app/

# Expose port 8080 to the outside world
EXPOSE 8080

# Copy the entrypoint script into the container
COPY entrypoint.sh .

# Make the entrypoint script executable
RUN chmod +x entrypoint.sh

# Set the entrypoint to wait for PostgreSQL and RabbitMQ and then run the executable
ENTRYPOINT ["./entrypoint.sh"]