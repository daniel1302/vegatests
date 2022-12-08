.PHONY: all
all:
	go mod tidy;
	go test -v ./...