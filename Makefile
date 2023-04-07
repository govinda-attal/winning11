SHELL:=bash
TAG_COMMIT := $(shell git rev-list --abbrev-commit --tags --max-count=1)
# `2>/dev/null` suppress errors and `|| true` suppress the error codes.
TAG := $(shell git describe --abbrev=0 --tags ${TAG_COMMIT} 2>/dev/null || true)
# here we strip the version prefix
VERSION := $(TAG:v%=%)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")

ifeq ($(VERSION),)
    VERSION := $(COMMIT)-$(DATE)
endif
ifneq ($(COMMIT), $(TAG_COMMIT))
    VERSION := $(VERSION)-next-$(COMMIT)-$(DATE)
endif
# git status --porcelain outputs a machine-readable text and the output is empty
# if the working tree is clean
ifneq ($(shell git status --porcelain),)
    VERSION := $(VERSION)-dirty
endif

FLAGS := -ldflags "-X github.com/govinda-attal/winning11/cmd.version=$(VERSION)"

tools:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.52.2
lint:
	@go vet ./...
	@golangci-lint run

test:
	@go test ./feeds/... ./utils/... -coverprofile cover.out
	@go tool cover -func cover.out | grep total:

build:
	@rm -rf dist
	@mkdir -p dist/configs
	@cp -r configs/*.yaml dist/configs
	@cp configs/*.json dist/
	@CGO_ENABLED=0 go build $(FLAGS) -o dist/winning11

docker_up:
	@docker compose up -d
docker_down:
	@docker compose down

run_local:
	@cd dist && APP_VALIDATOR_LOCAL=true ./winning11 validate --article valid-article.json
	@cd dist && APP_VALIDATOR_LOCAL=true ./winning11 validate --article invalid-article.json

run_db:
	@cd dist && ./winning11 migrate
	@cd dist && APP_VALIDATOR_LOCAL=false ./winning11 validate --article valid-article.json