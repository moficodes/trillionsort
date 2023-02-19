GENERATE := generate
JOINFILE := joinfile
SORTFILE := sortfile
SPLITFILE := splitfile
EXTERNALSORT := externalsort

clean:
	rm -f $(GENERATE)
	rm -f $(JOINFILE)
	rm -f $(SORTFILE)
	rm -f $(EXTERNALSORT)
	rm -f $(SPLITFILE)

generate:
	rm -f $(GENERATE)
	CGO_ENABLED=0 go build -o $(GENERATE) -ldflags="-s -w" cmd/generate/main.go

splitfile:
	rm -rf $(SPLITFILE)
	CGO_ENABLED=0 go build -o $(SPLITFILE) -ldflags="-s -w" cmd/splitfile/main.go

joinfile:
	rm -rf $(JOINFILE)
	CGO_ENABLED=0 go build -o $(JOINFILE) -ldflags="-s -w" cmd/joinfile/main.go

sortfile:
	rm -rf $(SORTFILE)
	CGO_ENABLED=0 go build -o $(SORTFILE) -ldflags="-s -w" cmd/sortfile/main.go

externalsort:
	rm -rf $(EXTERNALSORT)
	CGO_ENABLED=0 go build -o $(EXTERNALSORT) -ldflags="-s -w" cmd/externalsort/main.go

.PHONY: clean generate joinfile sortfile externalsort splitfile