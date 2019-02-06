SUBDIRS := ./lexer ./token ./ast ./repl ./parser
autotest:
	find . -iname '*.go' | entr -r go test -v --cover $(SUBDIRS)
