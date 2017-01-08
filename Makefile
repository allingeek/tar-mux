# needs Docker
PWD := $(shell pwd)

build:
	@docker run --rm \
	  -v $(PWD):/go/src/github.com/allingeek/tar-mux \
	  -v $(PWD)/bin:/go/bin \
	  -w /go/src/github.com/allingeek/tar-mux \
	  -e GOOS=linux \
	  -e GOARCH=amd64 \
	  golang:1.7 \
	  go build -ldflags="-s -w" -o bin/tar-mux-linux64
	@docker run --rm \
	  -v $(PWD):/go/src/github.com/allingeek/tar-mux \
	  -v $(PWD)/bin:/go/bin \
	  -w /go/src/github.com/allingeek/tar-mux \
	  -e GOOS=darwin \
	  -e GOARCH=amd64 \
	  golang:1.7 \
	  go build -ldflags="-s -w" -o bin/tar-mux-darwin64
upx: build
	@docker run --rm \
	  -v $(PWD)/bin:/input \
	  -w /input \
	  allingeek/upx:latest \
	  --brute -k tar-mux-linux64
	@mv ./bin/tar-mux-linux64 ./bin/tar-mux-linux64-upx
	@mv ./bin/tar-mux-linux64.~ ./bin/tar-mux-linux64
	@docker run --rm \
	  -v $(PWD)/bin:/input \
	  -w /input \
	  allingeek/upx:latest \
	  --brute -k tar-mux-darwin64
	@mv ./bin/tar-mux-darwin64 ./bin/tar-mux-darwin64-upx
	@mv ./bin/tar-mux-darwin64.~ ./bin/tar-mux-darwin64
