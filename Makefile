CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/sfomuseum/go-sfomuseum-flysfo
	cp *.go src/github.com/sfomuseum/go-sfomuseum-flysfo/
	cp -r client src/github.com/sfomuseum/go-sfomuseum-flysfo/
	cp -r archive src/github.com/sfomuseum/go-sfomuseum-flysfo/
	cp -r flight src/github.com/sfomuseum/go-sfomuseum-flysfo/
	cp -r iter src/github.com/sfomuseum/go-sfomuseum-flysfo/
	cp -r lookup src/github.com/sfomuseum/go-sfomuseum-flysfo/
	cp -r sfomuseum src/github.com/sfomuseum/go-sfomuseum-flysfo/
	cp -r storage src/github.com/sfomuseum/go-sfomuseum-flysfo/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/aaronland/go-storage"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-csv"
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/gjson"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt *.go

bin: 	self
	rm -rf bin/*
