BINARY      := seer-q
MODULE      := github.com/chris-xu0321/Midaz-cli
VERSION     := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
NPM_VERSION := $(shell node -p "require('./package.json').version" 2>/dev/null || echo dev)
LDFLAGS     := -s -w -X $(MODULE)/internal/build.Version=$(VERSION)
PREFIX      ?= /usr/local

.PHONY: build test clean release install qa

build:
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/seer-q/

test:
	go test -race -count=1 ./...

clean:
	rm -f $(BINARY) $(BINARY).exe
	rm -rf dist/ bin/

install: build
	mkdir -p $(PREFIX)/bin
	cp $(BINARY) $(PREFIX)/bin/$(BINARY)

release:
	GORELEASER_CURRENT_TAG=v$(NPM_VERSION) goreleaser release --clean --skip=publish

qa: test
	@echo "=== Skills validation test ==="
	bash test/skills-dist-test.sh
	@echo ""
	@echo "=== Smoke test (requires API) ==="
	bash test/smoke-test.sh ./$(BINARY) || echo "SKIP: API not running or binary not built"
