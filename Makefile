PHONY: brk prd cons

brk:
	go build -o ./bin/broker ./main.go

prd:
	go build -o ./bin/producer ./producer/producer.go

cons:
	go build -o ./bin/consumer ./consumer/consumer.go

all:
	make brk && make prd && make cons