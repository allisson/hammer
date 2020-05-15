build-protobuf:
	cd api/v1 && protoc -I/usr/local/include -I. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. hammer.proto
	cd api/v1 && protoc -I/usr/local/include -I. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. hammer.proto
	cd api/v1 && protoc -I/usr/local/include -I. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. hammer.proto

lint:
	if [ ! -f ./bin/golangci-lint ] ; \
	then \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.27.0; \
	fi;
	./bin/golangci-lint run

test:
	go test -covermode=count -coverprofile=count.out -v ./...

mock:
	@rm -rf mocks
	mockery -name TopicRepository
	mockery -name SubscriptionRepository
	mockery -name MessageRepository
	mockery -name DeliveryRepository
	mockery -name DeliveryAttemptRepository
	mockery -name TxRepository
	mockery -name TxFactoryRepository
	mockery -name TopicService
	mockery -name SubscriptionService
	mockery -name MessageService
	mockery -name DeliveryService
	mockery -name DeliveryAttemptService

run-worker:
	go run cmd/worker/main.go

run-server:
	go run cmd/server/main.go

.PHONY: build-protobuf lint test mock run-worker run-server
