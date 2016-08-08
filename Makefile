ON            = $(GOPATH)/bin/on
GO_BINDATA    = $(GOPATH)/bin/go-bindata
NODE_BIN      = $(shell npm bin)
PID_FILE      = .pid
GO_FILES      = $(filter-out ./bindata.go, $(shell find . -type f -name "*.go"))
TEMPLATES     = $(wildcard frontend/templates/*)
BINDATA       = bindata.go
BUNDLE        = build/bundle.js
APP           = $(wildcard frontend/*)
APP_NAME      = $(shell pwd | sed 's:.*/::')
TARGET        = ./kolide
GIT_HASH      = $(shell git rev-parse --short HEAD)
LDFLAGS       = -w -X main.commitHash=$(GIT_HASH)

build: $(ON) $(GO_BINDATA) clean $(TARGET)

clean:
	mkdir -p build
	rm -rf build/*
	rm -rf $(BINDATA)

$(ON):
	go get $(GOPATH)/src/vendor/github.com/olebedev/on

$(GO_BINDATA):
	go get $(GOPATH)/src/vendor/github.com/jteeuwen/go-bindata/...

$(BUNDLE): $(APP)
	$(NODE_BIN)/webpack --progress --colors --bail

$(TARGET): $(BUNDLE) $(BINDATA)
	go build -ldflags '$(LDFLAGS)' -o $@

kill:
	kill `cat $(PID_FILE)` || true

serve: $(ON) $(GO_BINDATA) clean $(BUNDLE) restart
	BABEL_ENV=dev node hot.proxy &
	$(NODE_BIN)/webpack --watch &
	$(ON) -m 2 $(GO_FILES) $(TEMPLATES) | xargs -n1 -I{} make restart || make kill

restart: $(BINDATA) kill $(TARGET)
	@echo restarting the app...
	$(TARGET) serve & echo $$! > $(PID_FILE)

$(BINDATA):
	$(GO_BINDATA) -o=$@ frontend/templates/ build/
