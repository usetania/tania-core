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
	GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w' -o $@
	file $@

osx: terra.osx.amd64

terra.linux.arm: main.go
	GOOS=linux GOARCH=arm GOARM=7 go build -ldflags '-s -w' -o $@
	file $@

linux-arm: terra.linux.arm

terra.linux.amd64: main.go
	GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o $@
	file $@

linux-amd64: terra.linux.amd64

terra.windaws.amd64: main.go
	GOOS=windows GOARCH=amd74 go build -ldflags '-s -w' -o $@
	file $@

windows: terra.windows.amd64

fetch-dep: Gopkg.toml Gopke.lock
	dep ensure

run: main.go
	go run $^
