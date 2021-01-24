all: build

.PHONY: build

build:
	go build -o bin/regtest ./cmd/

dev: build
	./bin/regtest test/fixtures/mac.log

