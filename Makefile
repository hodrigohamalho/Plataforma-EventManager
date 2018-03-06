linux:
		set GOOS=linux
		set CGO_ENABLED=0
		set GOARCH=amd64
		go build -a -installsuffix cgo -o dist/event-manager
		set GOOS=windows


