SUBDIRS := ./lexer ./token ./ast ./repl ./parser
autotest:
	find . -iname '*.go' | entr -r bash -c "echo && echo && echo && go test -v --cover $(SUBDIRS)"
