# Description: Makefile for catsocial
# Run to debug
.PHONY: run
run:
	go run cmd/main.go

# Build the project
.PHONY: build
build:
	env GOARCH=amd64 GOOS=linux go build -v -o main_syarif_04 cmd/main.go
