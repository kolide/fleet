.PHONY: build

export PATH := $(GOPATH)/bin:$(shell npm bin):$(PATH)

ifeq ($(OS), Windows_NT)
	OUTPUT = build/kolide.exe
else
	OUTPUT = build/kolide
endif

VERSION = 0.0.0-development
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
REVISION = $(shell git rev-parse HEAD)
USER = $(shell whoami)

ifeq ($(OS), Windows_NT)
	GOVERSION_CMD = "(go version).Split()[2]"
	GOVERSION = $(shell powershell $(GOVERSION_CMD))
	NOW	= $(shell powershell Get-Date -format s)
else
	GOVERSION = $(shell go version | awk '{print $$3}')
	NOW	= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
endif

DOCKER_IMAGE_NAME = kolide/kolide

ifndef CIRCLE_PR_NUMBER
	DOCKER_IMAGE_TAG = dev-unset
else
	DOCKER_IMAGE_TAG = dev-${CIRCLE_PR_NUMBER}
endif

all: build

define HELP_TEXT

  Makefile commands

	make deps         - Install depedent programs and libraries
	make generate     - Generate and bundle required code
	make generate-dev - Generate and bundle required code in a watch loop
	make distclean    - Delete all build artifacts

	make build        - Build the code

	make test         - Run the full test suite
	make test-go      - Run the Go tests
	make test-js      - Run the JavaScript tests
	make lint-go      - Run the Go linters
	make lint-js      - Run the JavaScript linters

endef

help:
	$(info $(HELP_TEXT))

.prefix:
ifeq ($(OS), Windows_NT)
	if not exist build mkdir build
else#
	mkdir -p build
endif

build: .prefix
	go build -i -o ${OUTPUT} -ldflags "\
	-X github.com/kolide/kolide-ose/version.version=${VERSION} \
	-X github.com/kolide/kolide-ose/version.branch=${BRANCH} \
	-X github.com/kolide/kolide-ose/version.revision=${REVISION} \
	-X github.com/kolide/kolide-ose/version.buildDate=${NOW} \
	-X github.com/kolide/kolide-ose/version.buildUser=${USER} \
	-X github.com/kolide/kolide-ose/version.goVersion=${GOVERSION}"

lint-js:
	eslint . --ext .js,.jsx

lint-go:
	go vet $(shell glide nv)

lint: lint-go lint-js

test-go:
	go test -v -cover $(shell glide nv)

test-js:
	_mocha --compilers js:babel-core/register \
		--recursive 'frontend/**/*.tests.js*' \
		--require 'frontend/.test.setup.js' \
		--require 'frontend/test/loaderMock.js'

test: lint test-go test-js

generate: .prefix
	go-bindata -pkg=server \
		-o=server/bindata.go \
		frontend/templates/ assets/...
	webpack --progress --colors --bail

generate-dev: .prefix
	go-bindata -debug -pkg=server \
		-o=server/bindata.go \
		frontend/templates/ assets/...
	webpack --progress --colors --bail --watch

deps:
	npm install
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/Masterminds/glide
	glide install

distclean:
ifeq ($(OS), Windows_NT)
	if exist build rmdir /s/q build
	if exist vendor rmdir /s/q vendor
	if exist assets\bundle.js del assets\bundle.js
else
	rm -rf build vendor
	rm -f assets/bundle.js
endif

docker-build-circle:
	@echo ">> building docker image"
	docker build -t "${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}" .
	docker push "${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}"