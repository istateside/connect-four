APP_NAME=connect_four

default: bin

bin:
	mkdir -p bin
	go build -i -o ./bin/connect_four

test:
	go test -timeout=10s -v

format:
	git ls-files | grep '.go$$' | xargs gofmt -w

deps:
	glide install

.PHONY: default bin test format deps

