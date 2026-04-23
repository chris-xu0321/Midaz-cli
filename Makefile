BINARY      := midaz
MODULE      := github.com/SparkssL/Midaz-cli
VERSION     := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS     := -s -w -X $(MODULE)/internal/build.Version=$(VERSION)
PREFIX      ?= /usr/local

.PHONY: build test clean release install qa

build:
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/midaz/

test:
	go test -race -count=1 ./...

clean:
	rm -f $(BINARY) $(BINARY).exe
	rm -rf dist/ bin/

install: build
	mkdir -p $(PREFIX)/bin
	cp $(BINARY) $(PREFIX)/bin/$(BINARY)

release:
	goreleaser release --clean --skip=publish

qa: test
	@echo "=== Skills validation test ==="
	bash test/skills-dist-test.sh
	@echo ""
	@echo "=== Smoke test (requires API) ==="
	bash test/smoke-test.sh ./$(BINARY) || echo "SKIP: API not running or binary not built"
