# Build Stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o books-api .

# Final Stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/books-api .
EXPOSE 8080
CMD ["./books-api"]
