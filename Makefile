LDFLAGS="-X github.com/g0dsCookie/gopbo/cmd.buildTime=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X github.com/g0dsCookie/gopbo/cmd.gitHash=`git rev-parse HEAD` -X github.com/g0dsCookie/gopbo/cmd.gitBranch=`git rev-parse --abbrev-ref HEAD`"

.PHONY: gopbo
gopbo: vendor
	go build -a -ldflags $(LDFLAGS) gopbo.go

.PHONY: all
all: gopbo-win32 gopbo-win64 gopbo-linux32 gopbo-linux64

.PHONY: gopbo-win32
gopbo-win32: vendor bindir
	env GOOS="windows" GOARCH="386" go build -a -ldflags $(LDFLAGS) -o gopbo.exe gopbo.go
	zip -9 bin/windows-x86.zip gopbo.exe
	rm -f gopbo.exe

.PHONY: gopbo-win64
gopbo-win64: vendor bindir
	env GOOS="windows" GOARCH="amd64" go build -a -ldflags $(LDFLAGS) -o gopbo.exe gopbo.go
	zip -9 bin/windows-x86_64.zip gopbo.exe
	rm -f gopbo.exe

.PHONY: gopbo-linux32
gopbo-linux32: vendor bindir
	env GOOS="linux" GOARCH="386" go build -a -ldflags $(LDFLAGS) -o gopbo gopbo.go
	zip -9 bin/linux-x86.zip gopbo
	rm -f gopbo

.PHONY: gopbo-linux64
gopbo-linux64: vendor bindir
	env GOOS="linux" GOARCH="amd64" go build -a -ldflags $(LDFLAGS) -o gopbo gopbo.go
	zip -9 bin/linux-x86_64.zip gopbo
	rm -f gopbo

.PHONY: bindir
bindir:
	mkdir -p bin/

.PHONY: vendor
vendor:
	git submodule update --init --recursive

.PHONY: clear
clear:
	rm -rf bin/