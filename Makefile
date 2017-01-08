# needs Docker
PWD := $(shell pwd)

build:
	@docker run --rm \
	  -v $(PWD):/go/src/github.com/allingeek/tar-stream-merge \
	  -v $(PWD)/bin:/go/bin \
	  -w /go/src/github.com/allingeek/tar-stream-merge \
	  -e GOOS=linux \
	  -e GOARCH=amd64 \
	  golang:1.7 \
	  go build -ldflags="-s -w" -o bin/tar-stream-merge-linux64
	@docker run --rm \
	  -v $(PWD):/go/src/github.com/allingeek/tar-stream-merge \
	  -v $(PWD)/bin:/go/bin \
	  -w /go/src/github.com/allingeek/tar-stream-merge \
	  -e GOOS=darwin \
	  -e GOARCH=amd64 \
	  golang:1.7 \
	  go build -ldflags="-s -w" -o bin/tar-stream-merge-darwin64
