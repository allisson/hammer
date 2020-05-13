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

We are using curl in the examples below.

### Create a new topic

```bash
curl -X POST 'http://localhost:8000/v1/topics' \
--header 'Content-Type: application/json' \
--data-raw '{
        "topic": {
                "id": "topic",
                "name": "Topic"
        }
}'
```

```javascript
{
  "id": "topic",
  "name": "Topic",
  "created_at": "2020-05-13T18:57:28.035492Z",
  "updated_at": "2020-05-13T18:57:28.035492Z"
}
```

### Create a new subscription

The max_delivery_attempts, delivery_attempt_delay and delivery_attempt_timeout are in seconds.

```bash
curl -X POST 'http://localhost:8000/v1/subscriptions' \
--header 'Content-Type: application/json' \
--data-raw '{
	"subscription": {
		"id": "httpbin-post",
		"topic_id": "topic",
		"name": "Httpbin Post",
		"url": "https://httpbin.org/post",
		"secret_token": "my-super-secret-token",
		"max_delivery_attempts": 5,
		"delivery_attempt_delay": 60,
		"delivery_attempt_timeout": 5
	}
}'
```

```javascript
{
  "id": "httpbin-post",
  "topic_id": "topic",
  "name": "Httpbin Post",
  "url": "https://httpbin.org/post",
  "secret_token": "my-super-secret-token",
  "max_delivery_attempts": 5,
  "delivery_attempt_delay": 60,
  "delivery_attempt_timeout": 5,
  "created_at": "2020-05-13T19:00:06.283020Z",
  "updated_at": "2020-05-13T19:00:06.283020Z"
}
```

### Create a new message

```bash
curl -X POST 'http://localhost:8000/v1/messages' \
--header 'Content-Type: application/json' \
--data-raw '{
	"message": {
		"topic_id": "topic",
		"data": "{\"name\": \"Allisson\"}"
	}
}'
```

```javascript
{
  "id": "01E87PJQERSRN3SXQTEZB5X4FF",
  "topic_id": "topic",
  "data": "eyJuYW1lIjogIkFsbGlzc29uIn0=",
  "created_at": "2020-05-13T19:01:03.064407Z"
}
```

###  Run the worker

The system will send a post request and the server must respond with the following status codes for the delivery to be considered successful: 200, 201, 202 and 204.

```bash
make run-worker
go run cmd/worker/main.go
{"level":"info","ts":1589396500.632363,"caller":"worker/main.go:179","msg":"worker-started"}
{"level":"info","ts":1589396501.42912,"caller":"worker/main.go:93","msg":"delivery-attempt-made","id":"01E87PJQET446MJ8G5BPN5VGGX","status":"completed","attempts":1,"max_delivery_attempts":5}
```

Submitted payload:

```javascript
{
  "created_at":"2020-05-13T16:01:03.066423-03:00",
  "data":"eyJuYW1lIjogIkFsbGlzc29uIn0=",
  "id":"01E87PKW55K0VX9F82NHE04VYR",
  "message_id":"01E87PJQERSRN3SXQTEZB5X4FF",
  "secret_token":"my-super-secret-token",
  "subscription_id":"httpbin-post",
  "topic_id":"topic"
}
```

### Get delivery data

```bash
curl -X GET http://localhost:8000/v1/deliveries/01E87PJQET446MJ8G5BPN5VGGX
```

```javascript
{
  "id": "01E87PJQET446MJ8G5BPN5VGGX",
  "topic_id": "topic",
  "subscription_id": "httpbin-post",
  "message_id": "01E87PJQERSRN3SXQTEZB5X4FF",
  "data": "eyJuYW1lIjogIkFsbGlzc29uIn0=",
  "url": "https://httpbin.org/post",
  "secret_token": "my-super-secret-token",
  "max_delivery_attempts": 5,
  "delivery_attempt_delay": 60,
  "delivery_attempt_timeout": 5,
  "scheduled_at": "2020-05-13T19:01:03.066423Z",
  "delivery_attempts": 1,
  "status": "completed",
  "created_at": "2020-05-13T19:01:03.066423Z",
  "updated_at": "2020-05-13T19:01:41.425554Z"
}
```

### Get delivery attempt data

The execution_duration are in milliseconds.

```bash
curl -X GET 'http://localhost:8000/v1/delivery-attempts?delivery_id=01E87PJQET446MJ8G5BPN5VGGX'
```

```javascript
{
  "delivery_attempts":[
    {
      "id":"01E87PKW55K0VX9F82NHE04VYR",
      "delivery_id":"01E87PJQET446MJ8G5BPN5VGGX",
      "request":"POST /post HTTP/1.1\r\nHost: httpbin.org\r\nContent-Type: application/json\r\n\r\n{\"id\":\"01E87PKW55K0VX9F82NHE04VYR\",\"topic_id\":\"topic\",\"subscription_id\":\"httpbin-post\",\"message_id\":\"01E87PJQERSRN3SXQTEZB5X4FF\",\"secret_token\":\"my-super-secret-token\",\"data\":\"eyJuYW1lIjogIkFsbGlzc29uIn0=\",\"created_at\":\"2020-05-13T16:01:03.066423-03:00\"}",
      "response":"HTTP/2.0 200 OK\r\nContent-Length: 990\r\nAccess-Control-Allow-Credentials: true\r\nAccess-Control-Allow-Origin: *\r\nContent-Type: application/json\r\nDate: Wed, 13 May 2020 19:01:41 GMT\r\nServer: gunicorn/19.9.0\r\n\r\n{\n  \"args\": {}, \n  \"data\": \"{\\\"id\\\":\\\"01E87PKW55K0VX9F82NHE04VYR\\\",\\\"topic_id\\\":\\\"topic\\\",\\\"subscription_id\\\":\\\"httpbin-post\\\",\\\"message_id\\\":\\\"01E87PJQERSRN3SXQTEZB5X4FF\\\",\\\"secret_token\\\":\\\"my-super-secret-token\\\",\\\"data\\\":\\\"eyJuYW1lIjogIkFsbGlzc29uIn0=\\\",\\\"created_at\\\":\\\"2020-05-13T16:01:03.066423-03:00\\\"}\", \n  \"files\": {}, \n  \"form\": {}, \n  \"headers\": {\n    \"Accept-Encoding\": \"gzip\", \n    \"Content-Length\": \"254\", \n    \"Content-Type\": \"application/json\", \n    \"Host\": \"httpbin.org\", \n    \"User-Agent\": \"Go-http-client/2.0\", \n    \"X-Amzn-Trace-Id\": \"Root=1-5ebc4415-ae6e2582cb39e4843440b382\"\n  }, \n  \"json\": {\n    \"created_at\": \"2020-05-13T16:01:03.066423-03:00\", \n    \"data\": \"eyJuYW1lIjogIkFsbGlzc29uIn0=\", \n    \"id\": \"01E87PKW55K0VX9F82NHE04VYR\", \n    \"message_id\": \"01E87PJQERSRN3SXQTEZB5X4FF\", \n    \"secret_token\": \"my-super-secret-token\", \n    \"subscription_id\": \"httpbin-post\", \n    \"topic_id\": \"topic\"\n  }, \n  \"origin\": \"177.37.153.46\", \n  \"url\": \"https://httpbin.org/post\"\n}\n",
      "response_status_code":200,
      "execution_duration":768,
      "success":true,
      "created_at":"2020-05-13T19:01:41.414943Z"
    }
  ]
}
```

## How to build docker images

```
docker build -f docker/server.Dockerfile -t hammer-server .
docker build -f docker/worker.Dockerfile -t hammer-worker .
docker build -f docker/migrate.Dockerfile -t hammer-migrate .
```
