runProducer:
	@pplog go run tgbot/cmd/main.go || true

runWorker:
	@pplog go run worker/cmd/main.go || true

redis:
	@docker run -d -p 6379:6379 --name redis-test-instance redis

kafka:
	@docker run -d -p 9092:9092 --name broker apache/kafka:latest

ginstall:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
