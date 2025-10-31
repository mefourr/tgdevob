runProducer:
	@pplog go run telegram-api/cmd/main.go || true

runWorker:
	@pplog go run worker/cmd/main.go || true

redis:
	@docker run -d -p 6379:6379 --name redis-test-instance redis

kafka:
	@docker run -d -p 9092:9092 --name broker apache/kafka:latest

ginstall:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

grpc:
	protoc proto/voice-msg-validator/v1/*.proto \
	 --go_out=voice-msg-validator/pb --go_opt=module=github.com/mefourr/tgdevob/msg/voice/validator/pb \
	 --go-grpc_out=voice-msg-validator/pb --go-grpc_opt=module=github.com/mefourr/tgdevob/msg/voice/validator/pb

#tt:
#	protoc validator/pb/*.proto \
#	 --go_out=test/ --go_opt=module=github.com/mefourr/tgdevob/validator/pb \
#	 --go-grpc_out=test/ --go-grpc_opt=module=github.com/mefourr/tgdevob/validator/pb