DESTDIR :=
prefix  := /usr/local
bindir  := ${prefix}/bin
mandir  := ${prefix}/share/man

.PHONY: build
build:
	GOOS= GOARCH= GOARM= GOFLAGS= go build -o bin/_awsd_prompt $<
	bash configure.sh ${DESTDIR}${bindir}

.PHONY: install
install: build
	install -d ${DESTDIR}${bindir}
	install -m755 bin/_awsd_prompt ${DESTDIR}${bindir}/

.PHONY: uninstall
uninstall:
	rm -f ${DESTDIR}${bindir}/_awsd
	rm -f ${DESTDIR}${bindir}/_awsd_prompt
