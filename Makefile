# Application meta
CURRENT_DIR := $(shell pwd)
APP         := $(shell basename ${CURRENT_DIR})

# Directories
PROTO_DIR        := proto
PROTO_EXTERN_DIR := $(PROTO_DIR)/extern
GEN_DIR          := genproto
SWAGGER_DIR      := doc/swagger

# Proto sources
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

# protoc include paths
PROTO_INCLUDE = \
	-I$(PROTO_DIR) \
	-I$(PROTO_EXTERN_DIR) \
	-I$(shell dirname $(shell which protoc))/../include

GO := go

-include .env

.PHONY: run
run:
	$(GO) run cmd/main.go

.PHONY: proto
proto: proto_extern
	@echo "==> Regenerating protos"
	rm -rf $(GEN_DIR)
	mkdir -p $(GEN_DIR) $(SWAGGER_DIR)
	protoc $(PROTO_INCLUDE) \
		--go_out=$(GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(GEN_DIR) --grpc-gateway_opt=paths=source_relative,generate_unbound_methods=true \
		--openapiv2_out=$(SWAGGER_DIR) \
		--openapiv2_opt=allow_merge=true,merge_file_name=swagger_docs,use_allof_for_refs=true,disable_service_tags=true,json_names_for_fields=false \
		$(PROTO_FILES)

.PHONY: swagger
swagger: proto_extern
	@echo "==> Regenerating swagger only"
	mkdir -p $(SWAGGER_DIR)
	protoc $(PROTO_INCLUDE) \
		--openapiv2_out=$(SWAGGER_DIR) \
		--openapiv2_opt=allow_merge=true,merge_file_name=swagger_docs \
		$(PROTO_FILES)

# Подтянуть внешние proto (google/api + protoc-gen-openapiv2)
.PHONY: proto_extern
proto_extern:
	@mkdir -p $(PROTO_EXTERN_DIR)/google/api
	@mkdir -p $(PROTO_EXTERN_DIR)/protoc-gen-openapiv2/options
	@test -f $(PROTO_EXTERN_DIR)/google/api/annotations.proto || curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto -o $(PROTO_EXTERN_DIR)/google/api/annotations.proto
	@test -f $(PROTO_EXTERN_DIR)/google/api/http.proto        || curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto        -o $(PROTO_EXTERN_DIR)/google/api/http.proto
	@test -f $(PROTO_EXTERN_DIR)/protoc-gen-openapiv2/options/annotations.proto || curl -sSL https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/v2/master/proto/protoc-gen-openapiv2/options/annotations.proto -o $(PROTO_EXTERN_DIR)/protoc-gen-openapiv2/options/annotations.proto
	@test -f $(PROTO_EXTERN_DIR)/protoc-gen-openapiv2/options/openapiv2.proto   || curl -sSL https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/v2/master/proto/protoc-gen-openapiv2/options/openapiv2.proto   -o $(PROTO_EXTERN_DIR)/protoc-gen-openapiv2/options/openapiv2.proto

.PHONY: clean
clean:
	rm -rf $(GEN_DIR) $(SWAGGER_DIR)

.PHONY: sqlc
sqlc:
	sqlc generate

.PHONY: migrate-up
migrate-up:
	migrate -path migrations -database $$POSTGRES_URL up

.PHONY: migrate-down
migrate-down:
	migrate -path migrations -database $$POSTGRES_URL down 1

.PHONY: tidy
tidy:
	$(GO) mod tidy