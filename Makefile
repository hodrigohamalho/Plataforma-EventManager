linux:
		set GOOS=linux
		go build -o dist/event-manager

disable_cgo:
		set CGO_ENABLED=0


