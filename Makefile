PKG := github.com/meowfaceman/conshim

test:
	go test -v $(PKG)/...

lint:
	golangci-lint run
