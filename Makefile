.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/user.proto

.PHONY: build
build: proto
	go build -o bin/server cmd/api/main.go

.PHONY: run
run: build
	./bin/server

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	rm -rf bin/ 