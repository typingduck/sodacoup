LINTER=~/go/bin/golangci-lint

default: test lint fmt build

build: sodacoup generator

sodacoup: sodacoup.go go* sodacouplib/*.go
	go build sodacoup.go

generator: generator.go go* sodacouplib/*.go
	go build generator.go

test: *.go go* sodacouplib/*.go
	go test ./sodacouplib/...

lint: *.go sodacouplib/*.go
	test -x ${LINTER} && \
		${LINTER} run generator.go && \
		${LINTER} run sodacoup.go && \
		${LINTER} run sodacouplib/... || \
		echo no linter

fmt:
	go fmt ./...

