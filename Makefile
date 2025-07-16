runConsumer:
	@pplog go run telegram-api/cmd/main.go

runWorker:
	@pplog go run worker/cmd/main.go