#### development stage
FROM golang:1.16-buster AS build-env

# set envvar
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GO111MODULE='on'

# set workdir
WORKDIR /code

# get project dependencies
COPY go.mod go.sum /code/
RUN go mod download

# copy files
COPY . /code

# generate binary
RUN go build -ldflags="-s -w" -o ./app ./cmd/hammer

# final stage
FROM gcr.io/distroless/base:nonroot
COPY --from=build-env /code/app /
COPY --from=build-env /code/db/migrations /db/migrations
# http server
EXPOSE 8000
# health check server
EXPOSE 9000
# metrics server
EXPOSE 4001
# grpc server
EXPOSE 50051
ENTRYPOINT ["/app"]
