bindir  := /usr/local/Cellar/awsd/0.0.2/bin

.PHONY: build
build:
	GOOS= GOARCH= GOARM= GOFLAGS= go build -o bin/_awsd_prompt $<
	bash configure.sh ${bindir}


.PHONY: install
install: build
	install -d ${bindir}
	install -m755 bin/_awsd_prompt ${bindir}/

.PHONY: uninstall
uninstall:
	rm -f ${bindir}/_awsd
	rm -f ${bindir}/_awsd_prompt
