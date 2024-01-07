SOURCE_PROTO_PATH = ../proto
GO_SERVER_PROTO_PATH = ./protogo
PROTO_NAME = goServer.proto
# PROTO_NAME = pyServer.proto

protos:
	protoc -I $(SOURCE_PROTO_PATH) \
	--go_out=$(GO_SERVER_PROTO_PATH) \
	--go_opt=paths=source_relative \
    --go-grpc_out=$(GO_SERVER_PROTO_PATH) \
	--go-grpc_opt=paths=source_relative \
	$(SOURCE_PROTO_PATH)/$(PROTO_NAME)

# protoc -I . --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative goServer.proto