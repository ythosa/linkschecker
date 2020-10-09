.PHONY: build
build:
	go build -o build/debug/linkschecker -v ./src/cmd/apiserver/main.go

.PHONY: windows-build
windows-build:
	env GOOS=windows GOARCH=amd64 go build -o build/release/linkschecker.exe -v ./src/cmd/apiserver/main.go

.PHONY: linux-build
linux-build:
	env GOOS=linux GOARCH=amd64 go build -o build/release/linkschecker -v ./src/cmd/apiserver/main.go

.PHONY: run
run:
	go run ./src/cmd/apiserver/main.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -v ./src/...

.PHONY: pipeline
pipeline:
	make lint && make test && make

.DEFAULT_GOAL := build
