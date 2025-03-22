# Use the official Golang image as the base image
FROM golang:1.24-alpine

# Set the working directory inside the container
WORKDIR /app

# Ensure go.mod exists and copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o books-api .

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./books-api"]