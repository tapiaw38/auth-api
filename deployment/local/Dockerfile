# Stage 1: Build the Go Application
FROM golang:1.19-alpine as builder

# Add Maintainer Info
LABEL maintainer="tapiaw38 Singh <tapiaw38@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY . /app/

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Install PostgreSQL client and Golang-Migrate
RUN apk update && apk add --no-cache postgresql-client curl && \
    curl -L -o /usr/local/bin/migrate https://github.com/golang-migrate/migrate/releases/download/v4.15.0/migrate.linux-amd64.tar.gz && \
    chmod +x /usr/local/bin/migrate

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the Go application
CMD ["go", "run", "cmd/main.go", "--host", "0.0.0.0"]