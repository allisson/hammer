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

We are using https://httpie.org/ in the examples below.

### Create a new topic

```bash
http POST http://localhost:8000/v1/topics topic:='{"id": "person", "name": "Person"}'
HTTP/1.1 200 OK
Content-Length: 117
Content-Type: application/json
Date: Tue, 12 May 2020 18:51:45 GMT
Grpc-Metadata-Content-Type: application/grpc

{
    "created_at": "2020-05-12T18:51:45.065849Z",
    "id": "person",
    "name": "Person",
    "updated_at": "2020-05-12T18:51:45.065849Z"
}
```

### Create a new subscription

The max_delivery_attempts, delivery_attempt_delay and delivery_attempt_timeout are in seconds.

```bash
http POST http://localhost:8000/v1/subscriptions subscription:='{"id": "httpbin-post", "topic_id": "person", "name": "Httpbin Post", "url": "https://httpbin.org/post", "secret_token": "my-super-secret-token", "max_delivery_attempts": 5, "delivery_attempt_delay": 60, "delivery_attempt_timeout": 5}'
HTTP/1.1 200 OK
Content-Length: 304
Content-Type: application/json
Date: Tue, 12 May 2020 18:54:22 GMT
Grpc-Metadata-Content-Type: application/grpc

{
    "created_at": "2020-05-12T18:54:22.788427Z",
    "delivery_attempt_delay": 60,
    "delivery_attempt_timeout": 5,
    "id": "httpbin-post",
    "max_delivery_attempts": 5,
    "name": "Httpbin Post",
    "secret_token": "my-super-secret-token",
    "topic_id": "person",
    "updated_at": "2020-05-12T18:54:22.788427Z",
    "url": "https://httpbin.org/post"
}
```

### Create a new message

```bash
http POST http://localhost:8000/v1/messages message:='{"topic_id": "person", "data": "{\"name\": \"Allisson\"}"}'
HTTP/1.1 200 OK
Content-Length: 136
Content-Type: application/json
Date: Tue, 12 May 2020 18:56:02 GMT
Grpc-Metadata-Content-Type: application/grpc

{
    "created_at": "2020-05-12T18:56:02.175614Z",
    "data": "eyJuYW1lIjogIkFsbGlzc29uIn0=",
    "id": "01E853WTKZVYWJFAQZ0EC4G125",
    "topic_id": "person"
}
```

###  Run the worker

The system will send a post request and the server must respond with the following status codes for the delivery to be considered successful: 200, 201, 202 and 204.

```bash
make run-worker
go run cmd/worker/main.go
{"level":"info","ts":1589309806.2644632,"caller":"worker/main.go:169","msg":"worker-started"}
{"level":"info","ts":1589309806.2704692,"caller":"worker/main.go:120","msg":"fetch_deliveries","count":1}
{"level":"info","ts":1589309806.857866,"caller":"worker/main.go:87","msg":"delivery-attempt-made","id":"01E853WTM362FWJVB41YREBFAS","status":"completed","attempts":1,"max_delivery_attempts":5}
{"level":"info","ts":1589309811.863363,"caller":"worker/main.go:120","msg":"fetch_deliveries","count":0}
{"level":"info","ts":1589309816.868973,"caller":"worker/main.go:120","msg":"fetch_deliveries","count":0}
```

### Get delivery data

```bash
http http://localhost:8000/v1/deliveries/01E853WTM362FWJVB41YREBFAS
HTTP/1.1 200 OK
Content-Length: 497
Content-Type: application/json
Date: Tue, 12 May 2020 18:57:58 GMT
Grpc-Metadata-Content-Type: application/grpc

{
    "created_at": "2020-05-12T18:56:02.179577Z",
    "data": "eyJuYW1lIjogIkFsbGlzc29uIn0=",
    "delivery_attempt_delay": 60,
    "delivery_attempt_timeout": 5,
    "delivery_attempts": 1,
    "id": "01E853WTM362FWJVB41YREBFAS",
    "max_delivery_attempts": 5,
    "message_id": "01E853WTKZVYWJFAQZ0EC4G125",
    "scheduled_at": "2020-05-12T18:56:02.179577Z",
    "secret_token": "my-super-secret-token",
    "status": "completed",
    "subscription_id": "httpbin-post",
    "topic_id": "person",
    "updated_at": "2020-05-12T18:56:46.854968Z",
    "url": "https://httpbin.org/post"
}
```

### Get delivery attempt data

The execution_duration are in milliseconds.

```bash
http http://localhost:8000/v1/delivery-attempts delivery_id==01E853WTM362FWJVB41YREBFAS
HTTP/1.1 200 OK
Content-Length: 1780
Content-Type: application/json
Date: Tue, 12 May 2020 19:00:22 GMT
Grpc-Metadata-Content-Type: application/grpc

{
    "delivery_attempts": [
        {
            "created_at": "2020-05-12T18:56:46.839077Z",
            "delivery_id": "01E853WTM362FWJVB41YREBFAS",
            "execution_duration": 566,
            "id": "01E853Y67Q70BTR4HGH1HSE4KV",
            "request": "POST /post HTTP/1.1\r\nHost: httpbin.org\r\nContent-Type: application/json\r\n\r\n{\"topic_id\":\"person\",\"subscription_id\":\"httpbin-post\",\"message_id\":\"01E853WTKZVYWJFAQZ0EC4G125\",\"secret_token\":\"\",\"data\":\"eyJuYW1lIjogIkFsbGlzc29uIn0=\",\"created_at\":\"2020-05-12T15:56:02.179577-03:00\"}",
            "response": "HTTP/2.0 200 OK\r\nContent-Length: 871\r\nAccess-Control-Allow-Credentials: true\r\nAccess-Control-Allow-Origin: *\r\nContent-Type: application/json\r\nDate: Tue, 12 May 2020 18:56:46 GMT\r\nServer: gunicorn/19.9.0\r\n\r\n{\n  \"args\": {}, \n  \"data\": \"{\\\"topic_id\\\":\\\"person\\\",\\\"subscription_id\\\":\\\"httpbin-post\\\",\\\"message_id\\\":\\\"01E853WTKZVYWJFAQZ0EC4G125\\\",\\\"secret_token\\\":\\\"\\\",\\\"data\\\":\\\"eyJuYW1lIjogIkFsbGlzc29uIn0=\\\",\\\"created_at\\\":\\\"2020-05-12T15:56:02.179577-03:00\\\"}\", \n  \"files\": {}, \n  \"form\": {}, \n  \"headers\": {\n    \"Accept-Encoding\": \"gzip\", \n    \"Content-Length\": \"200\", \n    \"Content-Type\": \"application/json\", \n    \"Host\": \"httpbin.org\", \n    \"User-Agent\": \"Go-http-client/2.0\", \n    \"X-Amzn-Trace-Id\": \"Root=1-5ebaf16e-ff0f55144a49777818244d5c\"\n  }, \n  \"json\": {\n    \"created_at\": \"2020-05-12T15:56:02.179577-03:00\", \n    \"data\": \"eyJuYW1lIjogIkFsbGlzc29uIn0=\", \n    \"message_id\": \"01E853WTKZVYWJFAQZ0EC4G125\", \n    \"secret_token\": \"\", \n    \"subscription_id\": \"httpbin-post\", \n    \"topic_id\": \"person\"\n  }, \n  \"origin\": \"177.37.153.46\", \n  \"url\": \"https://httpbin.org/post\"\n}\n",
            "response_status_code": 200,
            "success": true
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
