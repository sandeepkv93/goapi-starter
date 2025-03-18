# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git and dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies and tidy up
RUN go mod download && go mod tidy

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api cmd/api/main.go

# Final stage using distroless
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/bin/api .

# Copy the .env file
COPY .env .

# Use nonroot user (provided by distroless)
USER nonroot:nonroot

# Expose port
EXPOSE 3000

# Run the application
CMD ["./api"] 