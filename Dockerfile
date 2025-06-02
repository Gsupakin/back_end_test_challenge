FROM golang:1.24.3-alpine

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main ./cmd/api

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"] 