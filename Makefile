SERVER_DIR := server
GRPC_FILES := $(SERVER_DIR)/prattle.pb.go        \
	      $(SERVER_DIR)/prattle_grpc.pb.go

default: $(GRPC_FILES)

$(SERVER_DIR):
	mkdir -p $@

$(SERVER_DIR)/%.pb.go $(SERVER_DIR)/%_grpc.pb.go: protos/%.proto | $(SERVER_DIR)
	protoc -I protos/ $< --go_out=module=github.com/prattle-chat/prattle-proxy/server:$(SERVER_DIR) --go-grpc_out=module=github.com/prattle-chat/prattle-proxy/server:$(SERVER_DIR)
