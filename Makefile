test:
	go clean -testcache
	go test ./...

install:
	go build cmd/cli/main.go
	mv main /Users/christian/go/bin/zet