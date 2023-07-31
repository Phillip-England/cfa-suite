# Use the official Go image as the base image
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Copy the entire root directory into the container
COPY ./ /app

# Build the Go binary inside the /app directory
RUN go build -o main .

# Set the entry point to run the compiled binary
ENTRYPOINT ["./main"]