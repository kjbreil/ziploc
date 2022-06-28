DISTDIR=dist
RELEASEVERSION=$(shell echo ${VERSION})
DEBUGVERSION=$(shell git rev-list -1 HEAD | cut -c 1-10 )
RELEASEVERSION?=$(DEBUGVERSION)
HASH=$(shell git rev-list -1 HEAD)
GOVERSION=$(shell go version)
DATE=$(shell date)

build_ziploc_darwin:
	GOOS=darwin go build
build_ziploc_windows:
	GOOS=windows go build
build_ziploc_linux:
	GOOS=linux go build
dist_darwin: build_ziploc_darwin
	make distclean
dist_windows: build_ziploc_windows
	make distclean
dist_windows: build_ziploc_linux
	make distclean
distclean:
	find . ! -name 'ziploc' ! -name 'ziploc.exe' ! -name 'README.md' ! -name 'LICENSE' -type f -exec rm -f {} +
	find . ! -name '.git' -type d -exec rm -fR {} +
doc:
	godoc -http=localhost:6060
vet:
	go vet
