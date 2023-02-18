GENERATE := generate
JOINFILE := joinfile
SORTFILE := sortfile

clean:
	rm -f $(GENERATE)

generate:
	rm -f $(GENERATE)
	CGO_ENABLED=0 go build -o $(GENERATE) -ldflags="-s -w" cmd/generate/main.go

joinfile:
	rm -rf $(JOINFILE)
	CGO_ENABLED=0 go build -o $(JOINFILE) -ldflags="-s -w" cmd/joinfile/main.go

sortfile:
	rm -rf $(SORTFILE)
	CGO_ENABLED=0 go build -o $(SORTFILE) -ldflags="-s -w" cmd/sortfile/main.go

.PHONY: clean generate joinfile sortfile