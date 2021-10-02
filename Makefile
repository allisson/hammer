PLATFORM := $(shell uname | tr A-Z a-z)

build-protobuf:
	cd api/v1 && protoc -I/usr/local/include -I. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. hammer.proto
	cd api/v1 && protoc -I/usr/local/include -I. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. hammer.proto
	cd api/v1 && protoc -I/usr/local/include -I. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. hammer.proto

lint:
	if [ ! -f ./bin/golangci-lint ] ; \
	then \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.42.1; \
	fi;
	./bin/golangci-lint run

test:
	go test -covermode=count -coverprofile=count.out -v ./...

download-golang-migrate-binary:
	if [ ! -f ./migrate.$(PLATFORM)-amd64 ] ; \
	then \
		curl -sfL https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.$(PLATFORM)-amd64.tar.gz | tar -xvz; \
	fi;

db-migrate: download-golang-migrate-binary
	./migrate.$(PLATFORM)-amd64 -source file://db/migrations -database ${HAMMER_DATABASE_URL} up

db-test-migrate: download-golang-migrate-binary
	./migrate.$(PLATFORM)-amd64 -source file://db/migrations -database ${HAMMER_TEST_DATABASE_URL} up

mock:
	@rm -rf mocks
	mockery --all

run-worker:
	go run cmd/hammer/main.go worker

run-server:
	go run cmd/hammer/main.go server

run-migrate:
	go run cmd/hammer/main.go migrate

.PHONY: build-protobuf lint test download-golang-migrate-binary db-migrate db-test-migrate mock run-worker run-server run-migrate
