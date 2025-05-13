# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install Git for private module access
RUN apk add --no-cache git

# Copy go.mod and go.sum first
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Use build arg to specify which service to build
ARG SERVICE_PATH
RUN go build -o main ${SERVICE_PATH}

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
ENTRYPOINT ["./main"]
