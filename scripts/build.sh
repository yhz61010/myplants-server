#!/bin/bash

# Build script for MyPlants Server

echo "Building MyPlants Server..."
go build -o bin/myplants-server ./cmd/myplants-server

echo "Build complete."