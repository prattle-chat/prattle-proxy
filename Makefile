SERVER_DIR := server
GRPC_FILES := $(SERVER_DIR)/auth.pb.go             \
	      $(SERVER_DIR)/auth_grpc.pb.go        \
	      $(SERVER_DIR)/user.pb.go             \
	      $(SERVER_DIR)/user_grpc.pb.go        \
	      $(SERVER_DIR)/messaging.pb.go        \
	      $(SERVER_DIR)/messaging_grpc.pb.go   \
	      $(SERVER_DIR)/group.pb.go            \
	      $(SERVER_DIR)/group_grpc.pb.go

CERTS_DIR ?= certs
CERTS := $(CERTS_DIR)/ca-key.pem     \
	 $(CERTS_DIR)/ca-cert.pem    \
	 $(CERTS_DIR)/server-key.pem \
	 $(CERTS_DIR)/server-req.pem \
	 $(CERTS_DIR)/server-cert.pem

BINARY := prattle-proxy

default: $(GRPC_FILES) $(CERTS) $(BINARY)

$(SERVER_DIR):
	mkdir -p $@

$(SERVER_DIR)/%.pb.go $(SERVER_DIR)/%_grpc.pb.go: protos/%.proto | $(SERVER_DIR)
	protoc -I protos/ $< --go_out=module=github.com/prattle-chat/prattle-proxy/server:$(SERVER_DIR) --go-grpc_out=module=github.com/prattle-chat/prattle-proxy/server:$(SERVER_DIR)

$(CERTS_DIR) $(GENERATED_DIR):
	mkdir -p $@

$(CERTS): | $(CERTS_DIR)
	(cd $(CERTS_DIR) && ../scripts/gen-cert)


$(BINARY): $(GRPC_FILES) $(CERTS) *.go go.mod go.sum
	CGO_ENABLED=0 go build -ldflags="-s -w" -buildmode=pie -trimpath -o $@
