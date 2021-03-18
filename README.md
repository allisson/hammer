# hammer
[![Build Status](https://github.com/allisson/hammer/workflows/tests/badge.svg)](https://github.com/allisson/hammer/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/allisson/hammer)](https://goreportcard.com/report/github.com/allisson/hammer)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/allisson/hammer)
[![Docker Image](https://img.shields.io/docker/cloud/build/allisson/hammer)](https://hub.docker.com/r/allisson/hammer)

Simple webhook system written in golang.

## Features

- GRPC + Rest API.
- Topics and Subscriptions scheme similar to google pubsub, the message is published in a topic and this topic has several subscriptions sending the same notification to different systems.
- Payload sent follows the JSON Event Format for CloudEvents - Version 1.0 standard.
- Control the maximum amount of delivery attempts and delay between these attempts.
- Locks control of worker deliveries using PostgreSQL SELECT FOR UPDATE SKIP LOCKED.
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
docker run -p 4001:4001 -p 8000:8000 -p 9000:9000 -p 50051:50051 --env HAMMER_DATABASE_URL='postgres://user:pass@localhost:5432/hammer?sslmode=disable' allisson/hammer server # run grpc/http/metrics servers
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
  "created_at": "2021-03-18T11:08:49.678732Z",
  "updated_at": "2021-03-18T11:08:49.678732Z"
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
  "created_at": "2021-03-18T11:10:05.855296Z",
  "updated_at": "2021-03-18T11:10:05.855296Z"
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
  "id": "01F12GF6VAXGNHVXM4YT37N75A",
  "topic_id": "topic",
  "content_type": "application/json",
  "data": "eyJuYW1lIjogIkFsbGlzc29uIn0=",
  "created_at": "2021-03-18T11:10:29.738632Z"
}
```

###  Run the worker

The system will send a post request and the server must respond with the following status codes for the delivery to be considered successful: 200, 201, 202 and 204.

#### Docker

```bash
docker run --env HAMMER_DATABASE_URL='postgres://user:pass@localhost:5432/hammer?sslmode=disable' allisson/hammer worker
{"level":"info","ts":1616065862.332101,"caller":"service/worker.go:67","msg":"worker-started"}
{"level":"info","ts":1616065863.104438,"caller":"service/worker.go:36","msg":"worker-delivery-attempt-created","id":"01F12GG6NWZR03MW1MFMQDWVVF","delivery_id":"01F12GF6VM4YSX5GW8TM4781EZ","response_status_code":200,"execution_duration":749,"success":true}
```

#### Local

```bash
make run-worker
go run cmd/worker/main.go
{"level":"info","ts":1616065862.332101,"caller":"service/worker.go:67","msg":"worker-started"}
{"level":"info","ts":1616065863.104438,"caller":"service/worker.go:36","msg":"worker-delivery-attempt-created","id":"01F12GG6NWZR03MW1MFMQDWVVF","delivery_id":"01F12GF6VM4YSX5GW8TM4781EZ","response_status_code":200,"execution_duration":749,"success":true}
```

Submitted payload (Compatible with JSON Event Format for CloudEvents - Version 1.0):

```javascript
{
  "data_base64": "eyJuYW1lIjogIkFsbGlzc29uIn0=",
  "datacontenttype": "application/json",
  "id": "01F12GF6VM4YSX5GW8TM4781EZ",
  "messageid": "01F12GF6VAXGNHVXM4YT37N75A",
  "secrettoken": "my-super-secret-token",
  "source": "/v1/messages/01F12GF6VAXGNHVXM4YT37N75A",
  "specversion": "1.0",
  "subscriptionid": "httpbin-post",
  "time": "2021-03-18T11:10:29.748978Z",
  "topicid": "topic",
  "type": "hammer.message.created"
}
```

### Get delivery data

```bash
curl -X GET http://localhost:8000/v1/deliveries/01F12GF6VM4YSX5GW8TM4781EZ
```

```javascript
{
  "id": "01F12GF6VM4YSX5GW8TM4781EZ",
  "topic_id": "topic",
  "subscription_id": "httpbin-post",
  "message_id": "01F12GF6VAXGNHVXM4YT37N75A",
  "content_type": "application/json",
  "data": "eyJuYW1lIjogIkFsbGlzc29uIn0=",
  "url": "https://httpbin.org/post",
  "secret_token": "my-super-secret-token",
  "max_delivery_attempts": 5,
  "delivery_attempt_delay": 60,
  "delivery_attempt_timeout": 5,
  "scheduled_at": "2021-03-18T11:10:29.748978Z",
  "delivery_attempts": 1,
  "status": "completed",
  "created_at": "2021-03-18T11:10:29.748978Z",
  "updated_at": "2021-03-18T11:11:03.098648Z"
}
```

### Get delivery attempt data

The execution_duration are in milliseconds.

```bash
curl -X GET http://localhost:8000/v1/delivery-attempts/01F12GG6NWZR03MW1MFMQDWVVF
```

```javascript
{
  "id": "01F12GG6NWZR03MW1MFMQDWVVF",
  "delivery_id": "01F12GF6VM4YSX5GW8TM4781EZ",
  "request": "POST /post HTTP/1.1\r\nHost: httpbin.org\r\nContent-Type: application/json\r\n\r\n{\"specversion\":\"1.0\",\"type\":\"hammer.message.created\",\"source\":\"/v1/messages/01F12GF6VAXGNHVXM4YT37N75A\",\"id\":\"01F12GF6VM4YSX5GW8TM4781EZ\",\"time\":\"2021-03-18T11:10:29.748978Z\",\"secrettoken\":\"my-super-secret-token\",\"messageid\":\"01F12GF6VAXGNHVXM4YT37N75A\",\"subscriptionid\":\"httpbin-post\",\"topicid\":\"topic\",\"datacontenttype\":\"application/json\",\"data_base64\":\"eyJuYW1lIjogIkFsbGlzc29uIn0=\"}",
  "response": "HTTP/2.0 200 OK\r\nContent-Length: 1298\r\nAccess-Control-Allow-Credentials: true\r\nAccess-Control-Allow-Origin: *\r\nContent-Type: application/json\r\nDate: Thu, 18 Mar 2021 11:11:03 GMT\r\nServer: gunicorn/19.9.0\r\n\r\n{\n  \"args\": {}, \n  \"data\": \"{\\\"specversion\\\":\\\"1.0\\\",\\\"type\\\":\\\"hammer.message.created\\\",\\\"source\\\":\\\"/v1/messages/01F12GF6VAXGNHVXM4YT37N75A\\\",\\\"id\\\":\\\"01F12GF6VM4YSX5GW8TM4781EZ\\\",\\\"time\\\":\\\"2021-03-18T11:10:29.748978Z\\\",\\\"secrettoken\\\":\\\"my-super-secret-token\\\",\\\"messageid\\\":\\\"01F12GF6VAXGNHVXM4YT37N75A\\\",\\\"subscriptionid\\\":\\\"httpbin-post\\\",\\\"topicid\\\":\\\"topic\\\",\\\"datacontenttype\\\":\\\"application/json\\\",\\\"data_base64\\\":\\\"eyJuYW1lIjogIkFsbGlzc29uIn0=\\\"}\", \n  \"files\": {}, \n  \"form\": {}, \n  \"headers\": {\n    \"Accept-Encoding\": \"gzip\", \n    \"Content-Length\": \"386\", \n    \"Content-Type\": \"application/json\", \n    \"Host\": \"httpbin.org\", \n    \"User-Agent\": \"Go-http-client/2.0\", \n    \"X-Amzn-Trace-Id\": \"Root=1-60533547-501c866f62e44ea3736dbc0c\"\n  }, \n  \"json\": {\n    \"data_base64\": \"eyJuYW1lIjogIkFsbGlzc29uIn0=\", \n    \"datacontenttype\": \"application/json\", \n    \"id\": \"01F12GF6VM4YSX5GW8TM4781EZ\", \n    \"messageid\": \"01F12GF6VAXGNHVXM4YT37N75A\", \n    \"secrettoken\": \"my-super-secret-token\", \n    \"source\": \"/v1/messages/01F12GF6VAXGNHVXM4YT37N75A\", \n    \"specversion\": \"1.0\", \n    \"subscriptionid\": \"httpbin-post\", \n    \"time\": \"2021-03-18T11:10:29.748978Z\", \n    \"topicid\": \"topic\", \n    \"type\": \"hammer.message.created\"\n  }, \n  \"origin\": \"191.33.94.128\", \n  \"url\": \"https://httpbin.org/post\"\n}\n",
  "response_status_code": 200,
  "execution_duration": 749,
  "success": true,
  "created_at": "2021-03-18T11:11:03.091808Z"
}
```

## Environment variables

All environment variables is defined on file local.env.

## How to build docker image

```
docker build -f Dockerfile -t hammer .
```

## REST API

Default port: 8000

```bash
curl --location --request GET 'http://localhost:8000/v1/topics'
```

To disable the rest api, set the environment variable **HAMMER_REST_API_ENABLED** to false.

```bash
export HAMMER_REST_API_ENABLED='false'
```

## Prometheus metrics

Default port: 4001

```bash
curl --location --request GET 'http://localhost:4001/metrics'
```

To disable prometheus metrics, set the environment variable **HAMMER_METRICS_ENABLED** to false.

```bash
export HAMMER_METRICS_ENABLED='false'
```

## Health check

Default port: 9000

```bash
curl --location --request GET 'http://localhost:9000/liveness'
curl --location --request GET 'http://localhost:9000/readiness'
```

To disable health check, set the environment variable **HAMMER_HEALTH_CHECK_ENABLED** to false.

```bash
export HAMMER_HEALTH_CHECK_ENABLED='false'
```
