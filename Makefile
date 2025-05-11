.PHONY: run-consumer run-engine test

run-consumer:
	go run cmd/kafka-consumer/main.go

run-engine:
	go run cmd/engine/main.go

test:
	go test ./...

generate-data:
	go run scripts/generate_bids.go