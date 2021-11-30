BINDIR = /usr/local/bin

help:          ## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

install:       ## Install Target
	GOOS= GOARCH= GOARM= GOFLAGS= go build -o ${BINDIR}/_awsd_prompt
	cp scripts/_awsd ${BINDIR}/_awsd
	@echo " -=-=--=-=-=-=-=-=-=-=-=-=-=-=- "
	@echo "                                "
	@echo "   To Finish Installation add   "
	@echo "                                "
	@echo "  alias awsd=\"source _awsd\"   "
	@echo "                                "
	@echo " to your bash profile or zshrc  "
	@echo "   then open new terminal or    "
	@echo "       source that file         "
	@echo "                                "
	@echo " -=-=--=-=-=-=-=-=-=-=-=-=-=-=- "

uninstall:     ## Uninstall Target
	rm -f ${BINDIR}/_awsd
	rm -f ${BINDIR}/_awsd_prompt
