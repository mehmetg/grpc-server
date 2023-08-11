protobuf.out:
	PATH=$$PATH:$$HOME/go/bin/:/usr/local/opt/protobuf@3/bin/ \
	protoc --proto_path=./protobuf --go_out=./api --go_opt=paths=source_relative \
		--go-grpc_out=./api --go-grpc_opt=paths=source_relative \
		protobuf/grpc_server.proto