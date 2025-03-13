# Build stage
FROM golang:1.20-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o github-webhook-youtrack

# Final stage
FROM alpine:latest

# Add ca certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/github-webhook-youtrack .

# Copy action_mappings.json into the final image
COPY action_mappings.json .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./github-webhook-youtrack"]