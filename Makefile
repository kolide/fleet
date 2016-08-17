NPM_BIN = ./node_modules/.bin/
ifeq ($(OS), Windows_NT)
	NPM_BIN = node_modules\\.bin\\
endif

.prefix:
ifeq ($(OS), Windows_NT)
	if not exist build mkdir build
else
	mkdir -p build
endif

generate: .prefix
	$(NPM_BIN)webpack --progress --colors --bail
	go-bindata -pkg=app -o=app/bindata.go frontend/templates/ build/

deps:
	npm install
	go get -u github.com/tools/godep
	go get -u github.com/jteeuwen/go-bindata/...
ifneq ($(OS), Windows_NT)
	go get -u github.com/olebedev/on
endif

docs:
	$(NPM_BIN)jsdoc frontend -r -c .jsdoc.json -P package.json --verbose

lint:
	$(NPM_BIN)flow check
	$(NPM_BIN)eslint frontend

distclean:
	mkdir -p build
	rm -rf build/*

ifneq ($(OS), Windows_NT)

PID_FILE = build/kolide.pid
GO_FILES = $(filter-out ./bindata.go, $(shell find . -type f -name "*.go"))
TEMPLATES = $(wildcard frontend/templates/*)

stop:
	kill `cat $(PID_FILE)` || true

watch: .prefix
	BABEL_ENV=dev node tools/app/hot.proxy &
	$(WEBPACK) --watch &
	on -m 2 $(GO_FILES) $(TEMPLATES) | xargs -n1 -I{} make restart || make stop

restart: stop
	@echo restarting the app...
	kolide serve & echo $$! > $(PID_FILE)

endif