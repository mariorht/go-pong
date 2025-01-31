#!/bin/bash

# Create builds directory if it doesn't exist
mkdir -p builds

# Build the server
echo "Building server..."
go build -o builds/pong_server server/main.go

# Build the client
echo "Building client..."
go build -o builds/pong_client client/main.go

echo "Build completed."
