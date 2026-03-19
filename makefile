PROTO_DIR = "/usr/local/Cellar/protobuf/33.4_1/include"
GOOGLEAPIS_DIR = "/Users/andrefedev/Documents/Dev/googleapis"
APP_PROTO_DIR = "/Users/andrefedev/Documents/Development/muydelcampo/flutter/pkg/api/lib/src/api"

run:
	set -a; source .env; set +a; go run ./cmd/server

_gen_proto_v1:
	@rm -rf ./protobuf/gen
	@mkdir -p ./protobuf/def/v1
	@mkdir -p ./protobuf/gen/v1
	@PATH=/Users/andrefedev/go/bin:$$PATH protoc -I=./protobuf/def/v1 -I=$(GOOGLEAPIS_DIR) --go_out=./protobuf/gen/v1 --go_opt=paths=source_relative ./protobuf/def/v1/*.proto
	@PATH=/Users/andrefedev/go/bin:$$PATH protoc -I=./protobuf/def/v1 -I=$(GOOGLEAPIS_DIR) --go-grpc_out=./protobuf/gen/v1 --go-grpc_opt=paths=source_relative ./protobuf/def/v1/*.proto

_gen_proto_dart_v2:
	@rm -rf $(APP_PROTO_DIR)
	@mkdir -p $(APP_PROTO_DIR)
	@protoc -I=protobuf/def/v1 -I=$(PROTO_DIR) -I=$(GOOGLEAPIS_DIR) --dart_out=grpc:$(APP_PROTO_DIR) protobuf/def/v1/*.proto $(PROTO_DIR)/google/protobuf/*.proto $(GOOGLEAPIS_DIR)/google/type/*.proto
