GOMOD = go mod
GOBUILD = go build -v -trimpath -race
GOTEST = go test -v -cover -race

ROOT = $(shell git rev-parse --show-toplevel)
BIN = dist/gbrowse
CMD = "./cmd/gbrowse"

VERSION = $(shell git describe --tags --abbrev=0)
COMMIT = $(shell git rev-parse HEAD)
GOVERSION = $(shell go version)

LDFLAGS_PKG = main
LDFLAGS = -ldflags="-X '$(LDFLAGS_PKG).AuthorName=' -X '$(LDFLAGS_PKG).AuthorEmail=' -X '$(LDFLAGS_PKG).Version=$(VERSION)' -X '$(LDFLAGS_PKG).GoVersion=$(GOVERSION)' -X '$(LDFLAGS_PKG).Commit=$(COMMIT)' -X '$(LDFLAGS_PKG).Project=gbrowse' -X '$(LDFLAGS_PKG).GithubUser=berquerant'"

.PHONY: $(BIN)
$(BIN):
	$(GOBUILD) -o $@ $(LDFLAGS) $(CMD)

DOCKER_RUN = docker run --rm -v "$(ROOT)":/usr/src/myapp -w /usr/src/myapp
DOCKER_IMAGE = golang:1.20

.PHONY: test
test:
ifdef GO_DOCKER_TEST
	$(DOCKER_RUN) $(DOCKER_IMAGE) $(GOTEST) $(LDFLAGS) ./...
else
	$(GOTEST) $(LDFLAGS) ./...
endif

.PHONY: init
init:
	$(GOMOD) tidy

.PHONY: generate
generate:
	go generate ./...
