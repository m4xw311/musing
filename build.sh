#!/bin/bash

# Build script for Musings

echo "Building Musings application..."

# Build the application
go build -o musings cmd/musings/main.go

echo "Build complete! Run ./musings to use the application."