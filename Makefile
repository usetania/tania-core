.PHONY: all cover clean clean-osx clean-linux-arm clean-linux-amd64 clean-win64 \
	osx linux-amd64 linux-arm windows fetch-dep run

all: osx linux-amd64 linux-arm windows


clean:
	@[ -f terra.osx.amd64 ] && rm -f terra.osx.amd64 || true
	@[ -f terra.linux.arm ] && rm -f terra.linux.arm || true
	@[ -f terra.linux.amd64 ] && rm -f terra.linux.amd64 || true
	@[ -f terra.win.amd64.exe ] && rm -f terra.win.amd64.exe || true

clean-osx: terra.osx.amd64
	rm -rf $^

clean-linux-arm: terra.linux.arm
	rm -rf $^

clean-linux-amd64: terra.linux.amd64
	rm -rf $^

clean-win64: terra.win.amd64.exe
	rm -rf $^

terra.osx.amd64: main.go
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w' -o $@
	file $@

osx: terra.osx.amd64

terra.linux.arm: main.go
	CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 \
		CC=arm-linux-gnueabihf-gcc	\
		CXX=arm-linux-gnueabihf-g++ \
		go build -ldflags '-s -w' -o $@
	file $@

linux-arm: terra.linux.arm

terra.linux.amd64: main.go
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o $@
	file $@

linux-amd64: terra.linux.amd64

terra.windows.amd64.exe: main.go
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 \
		CC=x86_64-w64-mingw32-gcc \
		CXX=x86_64-w64-mingw32-g++ \
		go build -ldflags '-s -w' -o $@
	file $@

windows: terra.windows.amd64.exe

fetch-dep: Gopkg.toml Gopke.lock
	dep ensure

run: main.go
	go run $^
