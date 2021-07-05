all:
	go mod tidy
	go test ./...
	go fmt ./...
	go build
	./dpfctl completion zsh > completion.zsh
	./dpfctl completion bash > completion.bash