# hammer
[![Build Status](https://github.com/allisson/hammer/workflows/tests/badge.svg)](https://github.com/allisson/hammer/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/allisson/hammer)](https://goreportcard.com/report/github.com/allisson/hammer)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/allisson/hammer)
[![Docker Image](https://img.shields.io/docker/cloud/build/allisson/hammer)](https://hub.docker.com/r/allisson/hammer)

Simple webhook system written in golang.

## Features

- GRPC + GPRC Gateway (http2 or http1, you choose your destination).
- Topics and Subscriptions scheme similar to google pubsub, the message is published in a topic and this topic has several subscriptions sending the same notification to different systems.
- Payload sent follows the JSON Event Format for CloudEvents - Version 1.0 standard.
- Control the maximum amount of delivery attempts and delay between these attempts.
- Locks control of worker deliveries using https://github.com/allisson/go-pglock.
- Simplicity, it does the minimum necessary, it will not have authentication/permission scheme among other things, the idea is to use it internally in the cloud and not leave exposed.

## Quickstart

Let's start with the basic concepts, we have three main entities that we must know to start:

- Topic: A named resource to which messages are sent.
- Subscription: A named resource representing a subscription to specific topic.
- Message: The data that a publisher sends to a topic and is eventually delivered to subscribers.

### Run the server

To run the server it is necessary to have a database available from postgresql, in this example we will consider that we have a database called hammer running in localhost with user and password equal to user.

#### Docker

```bash
docker run --env HAMMER_DATABASE_URL='postgres://user:pass@localhost:5432/hammer?sslmode=disable' allisson/hammer migrate # run database migrations
docker run --env HAMMER_DATABASE_URL='postgres://user:pass@localhost:5432/hammer?sslmode=disable' allisson/hammer server # run grpc server
```

#### Local

```bash
git clone https://github.com/allisson/hammer
cd hammer
cp local.env .env # and edit .env
make run-migrate # create database schema
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
  "created_at": "2020-05-15T12:46:25.789143Z",
  "updated_at": "2020-05-15T12:46:25.789143Z"
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
  "created_at": "2020-05-15T12:46:46.469638Z",
  "updated_at": "2020-05-15T12:46:46.469638Z"
}
```

### Create a new message

```bash
curl -X POST 'http://localhost:8000/v1/messages' \
--header 'Content-Type: application/json' \
--data-raw '{
	"message": {
		"topic_id": "topic",
		"content_type": "application/json",
		"data": "{\"name\": \"Allisson\"}"
	}
}'
```

```javascript
{
  "id": "01E8C5ZKFFKHEAMSM0WRKCH7FQ",
  "topic_id": "topic",
  "content_type": "application/json",
  "data": "eyJuYW1lIjogIkFsbGlzc29uIn0=",
  "created_at": "2020-05-15T12:47:11.343464Z"
}
```

###  Run the worker

The system will send a post request and the server must respond with the following status codes for the delivery to be considered successful: 200, 201, 202 and 204.

#### Docker

```bash
docker run --env HAMMER_DATABASE_URL='postgres://user:pass@localhost:5432/hammer?sslmode=disable' allisson/hammer worker
```

#### Local

```bash
make run-worker
go run cmd/worker/main.go
{"level":"info","ts":1589546853.6263108,"caller":"worker/main.go:182","msg":"worker-started"}
{"level":"info","ts":1589546854.666185,"caller":"worker/main.go:93","msg":"delivery-made","id":"01E8C5ZKFHGFR4ZBZ17SGN4FAX","topic_id":"topic","subscription_id":"httpbin-post","message_id":"01E8C5ZKFFKHEAMSM0WRKCH7FQ","status":"completed","attempts":1,"max_delivery_attempts":5}
```

Submitted payload (Compatible with JSON Event Format for CloudEvents - Version 1.0):

```javascript
{
  "data_base64": "eyJuYW1lIjogIkFsbGlzc29uIn0=",
  "datacontenttype": "application/json",
  "id": "01E8C5ZKFHGFR4ZBZ17SGN4FAX",
  "messageid": "01E8C5ZKFFKHEAMSM0WRKCH7FQ",
  "secrettoken": "my-super-secret-token",
  "source": "/v1/messages/01E8C5ZKFFKHEAMSM0WRKCH7FQ",
  "specversion": "1.0",
  "subscriptionid": "httpbin-post",
  "time": "2020-05-15T09:47:11.345912-03:00",
  "topicid": "topic",
  "type": "hammer.message.created"
}
```

### Get delivery data

```bash
curl -X GET http://localhost:8000/v1/deliveries/01E8C5ZKFHGFR4ZBZ17SGN4FAX
```

```javascript
{
  "id": "01E8C5ZKFHGFR4ZBZ17SGN4FAX",
  "topic_id": "topic",
  "subscription_id": "httpbin-post",
  "message_id": "01E8C5ZKFFKHEAMSM0WRKCH7FQ",
  "content_type": "application/json",
  "data": "eyJuYW1lIjogIkFsbGlzc29uIn0=",
  "url": "https://httpbin.org/post",
  "secret_token": "my-super-secret-token",
  "max_delivery_attempts": 5,
  "delivery_attempt_delay": 60,
  "delivery_attempt_timeout": 5,
  "scheduled_at": "2020-05-15T12:47:11.345912Z",
  "delivery_attempts": 1,
  "status": "completed",
  "created_at": "2020-05-15T12:47:11.345912Z",
  "updated_at": "2020-05-15T12:47:34.662483Z"
}
```

### Get delivery attempt data

The execution_duration are in milliseconds.

```bash
curl -X GET 'http://localhost:8000/v1/delivery-attempts?delivery_id=01E8C5ZKFHGFR4ZBZ17SGN4FAX'
```

```javascript
{
  "delivery_attempts":[
    {
      "id": "01E8C60987C47VZ14RBA48M5EB",
      "delivery_id": "01E8C5ZKFHGFR4ZBZ17SGN4FAX",
      "request": "POST /post HTTP/1.1\r\nHost: httpbin.org\r\nContent-Type: application/json\r\n\r\n{\"specversion\":\"1.0\",\"type\":\"hammer.message.created\",\"source\":\"/v1/messages/01E8C5ZKFFKHEAMSM0WRKCH7FQ\",\"id\":\"01E8C5ZKFHGFR4ZBZ17SGN4FAX\",\"time\":\"2020-05-15T09:47:11.345912-03:00\",\"secrettoken\":\"my-super-secret-token\",\"messageid\":\"01E8C5ZKFFKHEAMSM0WRKCH7FQ\",\"subscriptionid\":\"httpbin-post\",\"topicid\":\"topic\",\"datacontenttype\":\"application/json\",\"data_base64\":\"eyJuYW1lIjogIkFsbGlzc29uIn0=\"}",
      "response": "HTTP/2.0 200 OK\r\nContent-Length: 1308\r\nAccess-Control-Allow-Credentials: true\r\nAccess-Control-Allow-Origin: *\r\nContent-Type: application/json\r\nDate: Fri, 15 May 2020 12:47:34 GMT\r\nServer: gunicorn/19.9.0\r\n\r\n{\n  \"args\": {}, \n  \"data\": \"{\\\"specversion\\\":\\\"1.0\\\",\\\"type\\\":\\\"hammer.message.created\\\",\\\"source\\\":\\\"/v1/messages/01E8C5ZKFFKHEAMSM0WRKCH7FQ\\\",\\\"id\\\":\\\"01E8C5ZKFHGFR4ZBZ17SGN4FAX\\\",\\\"time\\\":\\\"2020-05-15T09:47:11.345912-03:00\\\",\\\"secrettoken\\\":\\\"my-super-secret-token\\\",\\\"messageid\\\":\\\"01E8C5ZKFFKHEAMSM0WRKCH7FQ\\\",\\\"subscriptionid\\\":\\\"httpbin-post\\\",\\\"topicid\\\":\\\"topic\\\",\\\"datacontenttype\\\":\\\"application/json\\\",\\\"data_base64\\\":\\\"eyJuYW1lIjogIkFsbGlzc29uIn0=\\\"}\", \n  \"files\": {}, \n  \"form\": {}, \n  \"headers\": {\n    \"Accept-Encoding\": \"gzip\", \n    \"Content-Length\": \"391\", \n    \"Content-Type\": \"application/json\", \n    \"Host\": \"httpbin.org\", \n    \"User-Agent\": \"Go-http-client/2.0\", \n    \"X-Amzn-Trace-Id\": \"Root=1-5ebe8f66-7cd7a3a08196283065db25d0\"\n  }, \n  \"json\": {\n    \"data_base64\": \"eyJuYW1lIjogIkFsbGlzc29uIn0=\", \n    \"datacontenttype\": \"application/json\", \n    \"id\": \"01E8C5ZKFHGFR4ZBZ17SGN4FAX\", \n    \"messageid\": \"01E8C5ZKFFKHEAMSM0WRKCH7FQ\", \n    \"secrettoken\": \"my-super-secret-token\", \n    \"source\": \"/v1/messages/01E8C5ZKFFKHEAMSM0WRKCH7FQ\", \n    \"specversion\": \"1.0\", \n    \"subscriptionid\": \"httpbin-post\", \n    \"time\": \"2020-05-15T09:47:11.345912-03:00\", \n    \"topicid\": \"topic\", \n    \"type\": \"hammer.message.created\"\n  }, \n  \"origin\": \"177.37.153.46\", \n  \"url\": \"https://httpbin.org/post\"\n}\n",
      "response_status_code": 200,
      "execution_duration": 1010,
      "success": true,
      "created_at": "2020-05-15T12:47:34.651052Z"
    }
  ]
}
```

## How to build docker image

```
docker build -f Dockerfile -t hammer .
```
