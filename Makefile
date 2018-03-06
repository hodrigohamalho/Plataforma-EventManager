linux:
		set GOOS=linux
		set GCO_ENABLED=0
		go build -o dist/event-manager
		set GOOS=windows


