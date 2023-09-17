#!/usr/bin/env sh

mkdir -p bin
GOOS=linux   GOARCH=amd64 go build -o bin/gitar
GOOS=windows GOARCH=amd64 go build -o bin/gitar.exe
