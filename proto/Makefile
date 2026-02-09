validator:
	protoc -I voice-msg-validator \
		voice-msg-validator/v1/*.proto \
		--go_out=voice-msg-validator/pb --go_opt=paths=source_relative \
	 	--go-grpc_out=voice-msg-validator/pb --go-grpc_opt=paths=source_relative

auth_gen:
	protoc auth/v1/*.proto \
		--go_out=auth/pb --go_opt=paths=source_relative \
	 	--go-grpc_out=auth/pb --go-grpc_opt=paths=source_relative

audio_storage:
	protoc s3_storage/v1/*.proto \
		--go_out=s3_storage/pb --go_opt=paths=source_relative \
	 	--go-grpc_out=s3_storage/pb --go-grpc_opt=paths=source_relative
