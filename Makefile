tgbot:
	@pplog go run tgbot/cmd/main.go || true

worker:
	@pplog go run worker/cmd/main.go || true

validator:
	@pplog go run voice-msg-validator/cmd/main.go || true

s3:
	@pplog go run s3/cmd/main.go || true

dev:
	@osascript \
		-e 'tell application "Terminal"' \
		-e '  do script "cd $(CURDIR) && make tgbot"' \
		-e '  do script "cd $(CURDIR) && make validator"' \
		-e '  do script "cd $(CURDIR) && make worker"' \
		-e '  do script "cd $(CURDIR) && make s3"' \
		-e 'end tell'

redis:
	@docker run -d -p 6379:6379 --name redis-test-instance redis

kafka:
	@docker run -d -p 9092:9092 --name broker apache/kafka:latest

proto-install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

auth-build:
	docker build -f authentication/Dockerfile -t auth .

auth: auth-build
	docker run -d -t -i \
		-e AUTH_FILE=authorized_key.json \
		-e ID \
		-e SERVICE_ACCOUNT_ID \
		-p 5552:5552 \
		auth
