#!/bin/sh

echo "Building server..."
go build -o .build/server ./cmd/server/main.go
echo "DONE"

echo "Building db seed util..."
go build -o .build/dbseed ./cmd/dbseed/main.go
echo "DONE"
