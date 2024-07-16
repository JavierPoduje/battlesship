.DEFAULT_GOAL := build

.phony: fmt vet build
fmt:
	go fmt ./...
vet: fmt
	go vet ./...
build: vet
	go build -o cmd/main cmd/main.go
