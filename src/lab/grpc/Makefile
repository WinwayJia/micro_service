all: grpc-server grpc-client

grpc-server: server.go
	go build -o $@ $^

grpc-client: client.go
	go build -o $@ $^
