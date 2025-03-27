#!/bin/sh

CONFIG_DIR="$(pwd)"

export CONFIG_DIR

exec go test ./...
