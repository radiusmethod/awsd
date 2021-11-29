EXE =
ifeq ($(GOOS),windows)
EXE = .exe
endif

DESTDIR :=
prefix  := /Users/pjoyce/workspace/git/pjaudiomv/awsd
bindir  := ${prefix}/bin

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
