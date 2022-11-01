create:
    protoc --proto_path=proto proto/*.proto --go_out=gen/proto/
    protoc --proto_path=proto proto/*.proto --go-grpc_out=gen/proto/

clean:
    rn gen/proto/*.go