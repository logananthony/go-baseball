# Start from the official Go base image
FROM golang:1.24

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first (for caching dependencies)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go app
RUN go build -o main .

# Expose the port your app listens on
EXPOSE 8080

# Command to run the app
CMD ["./main"]
