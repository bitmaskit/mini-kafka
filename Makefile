PHONY: run

run:
	make build && ./bin/main

build:
	go build -o ./bin/main ./main.go

slow:
	curl http://localhost:8080/slow

veryslow:
	curl http://localhost:8080/veryslow
