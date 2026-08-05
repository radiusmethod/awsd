BINDIR = /usr/local/bin

help:          ## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

install:       ## Install Target
	GOOS= GOARCH= GOARM= GOFLAGS= go build -o ${BINDIR}/_awsd_prompt
	cp scripts/_awsd ${BINDIR}/_awsd
	cp scripts/_awsd_autocomplete ${BINDIR}/_awsd_autocomplete
	@echo " -=-=--=-=-=-=-=-=-=-=-=-=-=-=-=-=-=- "
	@echo "                                      "
	@echo "     To Finish Installation add       "
	@echo "                                      "
	@echo "   eval \"\$$(_awsd_prompt init zsh)\"    "
	@echo "                                      "
	@echo "   to your zshrc (or bash profile,    "
	@echo "   with 'init bash') then open a new  "
	@echo "   terminal or source that file       "
	@echo "                                      "
	@echo " -=-=--=-=-=-=-=-=-=-=-=-=-=-=-=-=-=- "

uninstall:     ## Uninstall Target
	rm -f ${BINDIR}/_awsd
	rm -f ${BINDIR}/_awsd_autocomplete
	rm -f ${BINDIR}/_awsd_prompt

.PHONY: test test-coverage docs
test:          ## Run tests
	go test ./...

docs:          ## Regenerate docs/*.md from the cobra command tree
	cd tools/gendocs && go run . ../../docs

test-coverage: ## Run tests with coverage report
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
