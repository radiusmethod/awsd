BINDIR = /usr/local/bin
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS = -s -w -X github.com/radiusmethod/awsd/src/cmd.version=$(VERSION)

help:          ## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

install:       ## Install Target
	GOOS= GOARCH= GOARM= GOFLAGS= go build -ldflags="$(LDFLAGS)" -o ${BINDIR}/awsd
	@echo " -=-=--=-=-=-=-=-=-=-=-=-=-=-=-=-=- "
	@echo "                                    "
	@echo "    To Finish Installation add      "
	@echo "                                    "
	@echo "    eval \"\$$(awsd init zsh)\"         "
	@echo "                                    "
	@echo "  to your zshrc (or bash profile,   "
	@echo "  with 'init bash') then open a new "
	@echo "  terminal or source that file      "
	@echo "                                    "
	@echo " -=-=--=-=-=-=-=-=-=-=-=-=-=-=-=-=- "

uninstall:     ## Uninstall Target
	rm -f ${BINDIR}/awsd

.PHONY: test test-coverage docs
test:          ## Run tests
	go test ./...

docs:          ## Regenerate docs/*.md from the cobra command tree
	cd tools/gendocs && go run . ../../docs

test-coverage: ## Run tests with coverage report
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
