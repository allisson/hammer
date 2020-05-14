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
  "created_at": "2020-05-14T00:25:32.385472Z",
  "updated_at": "2020-05-14T00:25:32.385472Z"
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
  "created_at": "2020-05-14T00:25:57.064271Z",
  "updated_at": "2020-05-14T00:25:57.064271Z"
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
  "id": "01E8896CAFMFA5S0HRZEPMG67V",
  "topic_id": "topic",
  "content_type": "application/json",
  "data": "eyJuYW1lIjogIkFsbGlzc29uIn0=",
  "created_at": "2020-05-14T00:26:21.391993Z"
}
```

###  Run the worker

The system will send a post request and the server must respond with the following status codes for the delivery to be considered successful: 200, 201, 202 and 204.

```bash
make run-worker
go run cmd/worker/main.go
{"level":"info","ts":1589416007.9288619,"caller":"worker/main.go:182","msg":"worker-started"}
{"level":"info","ts":1589416008.6146638,"caller":"worker/main.go:93","msg":"delivery-attempt-made","id":"01E8896CAJ5MW2DEW325EWJQ2M","topic_id":"topic","subscription_id":"httpbin-post","message_id":"01E8896CAFMFA5S0HRZEPMG67V","status":"completed","attempts":1,"max_delivery_attempts":5}
```

Submitted payload (Compatible with JSON Event Format for CloudEvents - Version 1.0):

```javascript
{
  "data_base64": "eyJuYW1lIjogIkFsbGlzc29uIn0=",
  "datacontenttype": "application/json",
  "id": "01E8897685Z6RJF138AHQX0B6H",
  "messageid": "01E8896CAFMFA5S0HRZEPMG67V",
  "secrettoken": "my-super-secret-token",
  "source": "/v1/deliveries/01E8896CAJ5MW2DEW325EWJQ2M",
  "specversion": "1.0",
  "subscriptionid": "httpbin-post",
  "time": "2020-05-13T21:26:21.394236-03:00",
  "topicid": "topic",
  "type": "hammer.deliveryAttempt.create"
}
```

### Get delivery data

```bash
curl -X GET http://localhost:8000/v1/deliveries/01E8896CAJ5MW2DEW325EWJQ2M
```

```javascript
{
  "id": "01E8896CAJ5MW2DEW325EWJQ2M",
  "topic_id": "topic",
  "subscription_id": "httpbin-post",
  "message_id": "01E8896CAFMFA5S0HRZEPMG67V",
  "content_type": "application/json",
  "data": "eyJuYW1lIjogIkFsbGlzc29uIn0=",
  "url": "https://httpbin.org/post",
  "secret_token": "my-super-secret-token",
  "max_delivery_attempts": 5,
  "delivery_attempt_delay": 60,
  "delivery_attempt_timeout": 5,
  "scheduled_at": "2020-05-14T00:26:21.394236Z",
  "delivery_attempts": 1,
  "status": "completed",
  "created_at": "2020-05-14T00:26:21.394236Z",
  "updated_at": "2020-05-14T00:26:48.610829Z"
}
```

### Get delivery attempt data

The execution_duration are in milliseconds.

```bash
curl -X GET 'http://localhost:8000/v1/delivery-attempts?delivery_id=01E8896CAJ5MW2DEW325EWJQ2M'
```

```javascript
{
  "delivery_attempts": [
    {
      "id": "01E8897685Z6RJF138AHQX0B6H",
      "delivery_id": "01E8896CAJ5MW2DEW325EWJQ2M",
      "request": "POST /post HTTP/1.1\r\nHost: httpbin.org\r\nContent-Type: application/json\r\n\r\n{\"specversion\":\"1.0\",\"type\":\"hammer.deliveryAttempt.create\",\"source\":\"/v1/deliveries/01E8896CAJ5MW2DEW325EWJQ2M\",\"id\":\"01E8897685Z6RJF138AHQX0B6H\",\"time\":\"2020-05-13T21:26:21.394236-03:00\",\"secrettoken\":\"my-super-secret-token\",\"messageid\":\"01E8896CAFMFA5S0HRZEPMG67V\",\"subscriptionid\":\"httpbin-post\",\"topicid\":\"topic\",\"datacontenttype\":\"application/json\",\"data_base64\":\"eyJuYW1lIjogIkFsbGlzc29uIn0=\"}",
      "response": "HTTP/2.0 200 OK\r\nContent-Length: 1326\r\nAccess-Control-Allow-Credentials: true\r\nAccess-Control-Allow-Origin: *\r\nContent-Type: application/json\r\nDate: Thu, 14 May 2020 00:26:48 GMT\r\nServer: gunicorn/19.9.0\r\n\r\n{\n  \"args\": {}, \n  \"data\": \"{\\\"specversion\\\":\\\"1.0\\\",\\\"type\\\":\\\"hammer.deliveryAttempt.create\\\",\\\"source\\\":\\\"/v1/deliveries/01E8896CAJ5MW2DEW325EWJQ2M\\\",\\\"id\\\":\\\"01E8897685Z6RJF138AHQX0B6H\\\",\\\"time\\\":\\\"2020-05-13T21:26:21.394236-03:00\\\",\\\"secrettoken\\\":\\\"my-super-secret-token\\\",\\\"messageid\\\":\\\"01E8896CAFMFA5S0HRZEPMG67V\\\",\\\"subscriptionid\\\":\\\"httpbin-post\\\",\\\"topicid\\\":\\\"topic\\\",\\\"datacontenttype\\\":\\\"application/json\\\",\\\"data_base64\\\":\\\"eyJuYW1lIjogIkFsbGlzc29uIn0=\\\"}\", \n  \"files\": {}, \n  \"form\": {}, \n  \"headers\": {\n    \"Accept-Encoding\": \"gzip\", \n    \"Content-Length\": \"400\", \n    \"Content-Type\": \"application/json\", \n    \"Host\": \"httpbin.org\", \n    \"User-Agent\": \"Go-http-client/2.0\", \n    \"X-Amzn-Trace-Id\": \"Root=1-5ebc9048-8321b544c9aeb5781cc071ca\"\n  }, \n  \"json\": {\n    \"data_base64\": \"eyJuYW1lIjogIkFsbGlzc29uIn0=\", \n    \"datacontenttype\": \"application/json\", \n    \"id\": \"01E8897685Z6RJF138AHQX0B6H\", \n    \"messageid\": \"01E8896CAFMFA5S0HRZEPMG67V\", \n    \"secrettoken\": \"my-super-secret-token\", \n    \"source\": \"/v1/deliveries/01E8896CAJ5MW2DEW325EWJQ2M\", \n    \"specversion\": \"1.0\", \n    \"subscriptionid\": \"httpbin-post\", \n    \"time\": \"2020-05-13T21:26:21.394236-03:00\", \n    \"topicid\": \"topic\", \n    \"type\": \"hammer.deliveryAttempt.create\"\n  }, \n  \"origin\": \"177.37.153.46\", \n  \"url\": \"https://httpbin.org/post\"\n}\n",
      "response_status_code": 200,
      "execution_duration": 657,
      "success": true,
      "created_at": "2020-05-14T00:26:48.599343Z"
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
