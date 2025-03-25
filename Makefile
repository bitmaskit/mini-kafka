PHONY: run

run:
	make build && ./bin/main

build:
	go build -o ./bin/main ./main.go