SHELL = /bin/bash
NODE_BINDIR = ./node_modules/.bin
export PATH := $(NODE_BINDIR):$(PATH)
LOGNAME ?= $(shell logname)

.PHONY: all cover clean clean-osx clean-linux-arm clean-linux-amd64 clean-win64 \
	osx linux-amd64 linux-arm windows fetch-dep run osxcross.bin \
  cleantranslations makemessages translations

all: osx linux-amd64 linux-arm windows

clean:
	@[ -f tania.osx.amd64 ] && rm -f tania.osx.amd64 || true
	@[ -f tania.linux.arm ] && rm -f tania.linux.arm || true
	@[ -f tania.linux.amd64 ] && rm -f tania.linux.amd64 || true
	@[ -f tania.win.amd64.exe ] && rm -f tania.win.amd64.exe || true

clean-osx: tania.osx.amd64
	rm -rf $^

clean-linux-arm: tania.linux.arm
	rm -rf $^

clean-linux-amd64: tania.linux.amd64
	rm -rf $^

clean-win64: tania.win.amd64.exe
	rm -rf $^

tania.osx.amd64: main.go
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w' -o $@
	file $@

osx: tania.osx.amd64

osxcross: main.go
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 \
		CC=o64-clang	\
		CXX=o64-clang++ \
		go build -ldflags '-s -w' -o tania.osx.amd64
	file tania.osx.amd64

tania.linux.arm: main.go
	CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 \
		CC=arm-linux-gnueabihf-gcc	\
		CXX=arm-linux-gnueabihf-g++ \
		go build -ldflags '-s -w' -o $@
	file $@

linux-arm: tania.linux.arm

tania.linux.amd64: main.go
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o $@
	file $@

linux-amd64: tania.linux.amd64

tania.windows.amd64.exe: main.go
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 \
		CC=x86_64-w64-mingw32-gcc \
		CXX=x86_64-w64-mingw32-g++ \
		go build -ldflags '-s -w' -o $@
	file $@

windows: tania.windows.amd64.exe

fetch-dep: Gopkg.toml Gopke.lock
	dep ensure

run: main.go
	go run $^
