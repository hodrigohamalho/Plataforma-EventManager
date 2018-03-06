linux:
		set GOOS=linux
		set CGO_ENABLED=0
		set GOARCH=amd64
		go build -o dist/event-manager -v
		set GOOS=windows


