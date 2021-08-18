all: test fmt
	go mod tidy
	go build
	./dpfctl completion zsh > completion.zsh
	./dpfctl completion bash > completion.bash

test:
	go test -coverprofile="cover.out" ./...
	go tool cover -html=cover.out -o cover.html

fmt:
	go fmt ./...
