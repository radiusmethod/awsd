BINDIR  = /usr/local/bin

.PHONY: build
build:
	GOOS= GOARCH= GOARM= GOFLAGS= go build -o ${BINDIR}/_awsd_prompt


.PHONY: install
install: build
	chmod 755 ${BINDIR}/_awsd_prompt
	cp scripts/_awsd ${BINDIR}/_awsd
	@echo 'alias awsd="source _awsd"' >> ${HOME}/.zshrc
	@echo 'alias awsd="source _awsd"' >> ${HOME}/.bashrc
	@echo 'alias awsd="source _awsd"' >> ${HOME}/.bash_profile

.PHONY: uninstall
uninstall:
	rm -f ${BINDIR}/_awsd
	rm -f ${BINDIR}/_awsd_prompt
