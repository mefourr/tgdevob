runProducer:
	@pplog go run telegram-api/cmd/main.go

runWorker:
	@pplog go run worker/cmd/main.go

redis:
	@docker run -d -p 6379:6379 --name redis-test-instance redis

kafka:
	@docker run -d -p 9092:9092 --name broker apache/kafka:latest