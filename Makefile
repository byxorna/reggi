git_commit := $(shell git rev-parse HEAD)
git_tag := $(shell git describe --tags --always HEAD)
date := $(shell date)
pkg := $(shell go list -m)

all: build

.PHONY: build

build:
	@go build -o bin/reggi \
		-ldflags "-X '$(pkg)/pkg/version.Commit=$(git_commit)' -X '$(pkg)/pkg/version.Date=$(date)' -X '$(pkg)/pkg/version.Version=$(git_tag)'" ./cmd/reggi

dev: build
	@./bin/reggi pkg/ui/model.go test/fixtures/mac.log test/fixtures/test.txt

