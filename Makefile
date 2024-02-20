proto-gen:
	protoc --proto_path=protobufs  --go-grpc_out=internal/grpc --go_out=internal/grpc protobufs/account.proto

run:
	go run main.go