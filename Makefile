PLATFORM := $(shell uname | tr A-Z a-z)

lint:
	if [ ! -f ./bin/golangci-lint ] ; \
	then \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.26.0; \
	fi;
	./bin/golangci-lint run

test:
	go test -covermode=count -coverprofile=count.out -v ./...

download-golang-migrate-binary:
	if [ ! -f ./migrate.$(PLATFORM)-amd64 ] ; \
	then \
		curl -sfL https://github.com/golang-migrate/migrate/releases/download/v4.11.0/migrate.$(PLATFORM)-amd64.tar.gz | tar -xvz; \
	fi;

db-migrate: download-golang-migrate-binary
	./migrate.$(PLATFORM)-amd64 -source file://db/migrations -database ${DATABASE_URL} up

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

run:
	go run cmd/hammer/main.go

.PHONY: lint test download-golang-migrate-binary db-migrate mock run
