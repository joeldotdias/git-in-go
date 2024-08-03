cli:build
	@cp ./bin/gat ~/bin/projects/gat

build:
	@go build -o ./bin/gat cmd/gat/main.go
