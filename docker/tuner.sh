#!/bin/bash

apt update

apt-get clean && rm -rf /var/lib/apt/lists/*

go mod tidy

mkdir -p bin

go build -o ./bin/app ./cmd/main.go
