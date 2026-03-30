BINARY      := seer-q
MODULE      := github.com/chris-xu0321/Midaz-cli
VERSION     := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
NPM_VERSION := $(shell node -p "require('./npm/package.json').version" 2>/dev/null || echo dev)
LDFLAGS     := -s -w -X $(MODULE)/internal/build.Version=$(VERSION)

.PHONY: build test clean release qa qa-release

build:
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/seer-q/

test:
	go test -race -count=1 ./...

clean:
	rm -f $(BINARY) $(BINARY).exe
	rm -rf dist/

release:
	GORELEASER_CURRENT_TAG=v$(NPM_VERSION) goreleaser release --clean --skip=publish

qa: test
	@echo "=== Skills distribution test ==="
	bash test/skills-dist-test.sh
	@echo ""
	@echo "=== Smoke test (requires API) ==="
	bash test/smoke-test.sh ./$(BINARY) || echo "SKIP: API not running or binary not built"

qa-release: release qa
	@echo "=== npm package verification ==="
	bash npm/verify.sh
	@echo ""
	@echo "=== npm install test ==="
	bash test/npm-install-test.sh
