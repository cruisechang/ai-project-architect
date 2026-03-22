.PHONY: help test build build-mac build-linux build-windows build-all clean

TARGET ?= mac

help:
	@echo "Available targets:"
	@echo "  make test"
	@echo "  make build TARGET=mac"
	@echo "  make build TARGET=linux"
	@echo "  make build TARGET=windows"
	@echo "  make build-all"
	@echo "  make clean"

test:
	go test ./...
	./scripts/build_test.sh

build:
	./build.sh $(TARGET)

build-mac:
	./build.sh mac

build-linux:
	./build.sh linux

build-windows:
	./build.sh windows

build-all:
	./build.sh all

clean:
	rm -rf release build
