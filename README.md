# hammer
[![Build Status](https://github.com/allisson/hammer/workflows/tests/badge.svg)](https://github.com/allisson/hammer/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/allisson/hammer)](https://goreportcard.com/report/github.com/allisson/hammer)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/allisson/hammer)

Simple webhook system written in golang.

## Quickstart

Let's start with the basic concepts, we have three main entities that we must know to start:

- Topic: A named resource to which messages are sent.
- Subscription: A named resource representing a subscription to specific topic.
- Message: The data that a publisher sends to a topic and is eventually delivered to subscribers.

### Run the server

To run the server it is necessary to have a database available from postgresql, in this example we will consider that we have a database called hammer running in localhost with user and password equal to user.

```bash
git clone https://github.com/allisson/hammer
cd hammer
cp local.env .env # and edit .env
export HAMMER_DATABASE_URL='postgres://user:password@localhost:5432/hammer?sslmode=disable'
make db-migrate # create database schema
make run-server # run the server (grpc + http)
```

### Create a new topic

```bash
curl --location --request POST 'http://localhost:8000/v1/topics' \
--header 'Content-Type: application/json' \
--data-raw '{
	"topic": {
		"id": "person",
		"name": "Person"
	}
}'

{"id":"person","name":"Person","created_at":"2020-05-09T19:27:18.115663Z","updated_at":"2020-05-09T19:27:18.115663Z"}
```

### Create a new subscription

The max_delivery_attempts, delivery_attempt_delay and delivery_attempt_timeout are in seconds.

```bash
curl --location --request POST 'http://localhost:8000/v1/subscriptions' \
--header 'Content-Type: application/json' \
--data-raw '{
        "subscription": {
                "id": "httpbin-post",
                "topic_id": "person",
                "name": "Httpbin Post",
                "url": "https://httpbin.org/post",
                "secret_token": "my-super-secret-token",
                "max_delivery_attempts": 5,
                "delivery_attempt_delay": 60,
                "delivery_attempt_timeout": 5
        }
}'

{"id":"httpbin-post","topic_id":"person","name":"Httpbin Post","url":"https://httpbin.org/post","secret_token":"my-super-secret-token","max_delivery_attempts":5,"delivery_attempt_delay":60,"delivery_attempt_timeout":5,"created_at":"2020-05-09T20:15:30.057053Z","updated_at":"2020-05-09T20:15:30.057053Z"}
```

### Create a new message

```bash
curl --location --request POST 'http://localhost:8000/v1/messages' \
--header 'Content-Type: application/json' \
--data-raw '{
        "message": {
                "topic_id": "person",
                "data": "{\"name\": \"Allisson\"}"
        }
}'

{"id":"01E7XHBG5AYVWPYST6PZXZ2MRE","topic_id":"person","data":"{\"name\": \"Allisson\"}","created_at":"2020-05-09T20:17:19.018923Z"}
```

###  Run the worker

```bash
make run-worker
go run cmd/worker/main.go
{"level":"info","ts":1589055518.551332,"caller":"worker/main.go:179","msg":"worker-started"}
{"level":"info","ts":1589055518.555826,"caller":"worker/main.go:127","msg":"fetch_deliveries","count":1}
{"level":"info","ts":1589055519.437763,"caller":"worker/main.go:94","msg":"delivery-attempt-made","id":"01E7XHBG5DQTFH9KV3395TF096","status":"completed","attempts":1,"max_delivery_attempts":5}
{"level":"info","ts":1589055524.446191,"caller":"worker/main.go:127","msg":"fetch_deliveries","count":0}
{"level":"info","ts":1589055529.448663,"caller":"worker/main.go:127","msg":"fetch_deliveries","count":0}
```
