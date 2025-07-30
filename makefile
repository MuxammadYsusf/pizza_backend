proto_gen:
	protoc \
  --proto_path=proto --go_out=generated --go_opt=paths=source_relative \
  --go-grpc_out=generated --go-grpc_opt=paths=source_relative \
   proto/**/*.proto
	protoc \
  --proto_path=proto --go_out=generated --go_opt=paths=source_relative \
  --go-grpc_out=generated --go-grpc_opt=paths=source_relative \
   proto/**/*.proto