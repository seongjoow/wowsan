# Broker Proto
SOURCE_PROTO_PATH = ./pkg/proto/broker
GO_SERVER_PROTO_PATH = ./pkg/proto/broker
GO_PROTO_NAME = broker.proto

# # Publisher Proto
# SOURCE_PROTO_PATH = ./pkg/proto/publisher
# GO_SERVER_PROTO_PATH = ./pkg/proto/publisher
# GO_PROTO_NAME = publisher.proto

# # # Subscriber Proto
# SOURCE_PROTO_PATH = ./pkg/proto/subscriber
# GO_SERVER_PROTO_PATH = ./pkg/proto/subscriber
# GO_PROTO_NAME = subscriber.proto


# SOURCE_PROTO_PATH = ../proto
# GO_SERVER_PROTO_PATH = ./proto
# PROTO_NAME = goServer.proto
# # PROTO_NAME = pyServer.proto

protos:
	protoc -I $(SOURCE_PROTO_PATH) \
	--go_out=$(GO_SERVER_PROTO_PATH) \
	--go_opt=paths=source_relative \
    --go-grpc_out=$(GO_SERVER_PROTO_PATH) \
	--go-grpc_opt=paths=source_relative \
	$(SOURCE_PROTO_PATH)/$(GO_PROTO_NAME)