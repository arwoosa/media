# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/media ./main.go

# Stage 2: Create the final image
FROM scratch

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/media .

# Copy the proto files
COPY proto /app/proto

# Expose the port the application runs on
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/app/media"]
