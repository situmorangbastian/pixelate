# Use an official Golang Alpine image as a parent image
FROM golang:alpine

# Install FFmpeg and other necessary packages
RUN apk update && \
    apk add --no-cache ffmpeg && \
    rm -rf /var/cache/apk/*


# Set the working directory in the container
WORKDIR /go/src/app

# Copy the local package files to the container's workspace
COPY . .

# Rename config.toml.example to config.toml
RUN mv config.toml.example config.toml

# Build the Go application
RUN go build -o pixelate app/main.go

# Run the Go application
CMD ["./pixelate"]

