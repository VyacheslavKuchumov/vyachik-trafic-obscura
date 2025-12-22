build client:
	@go build -o bin/client ./cmd/client

build server:
	@go build -o bin/server ./cmd/server

run client: build client
	@./bin/client
	
run server: build server
	@./bin/server