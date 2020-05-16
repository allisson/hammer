FROM golang:1.14-buster AS build-env

# set envvar
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GO111MODULE='on'

# set workdir
WORKDIR /build

# copy files
COPY . .

# build binary
RUN go mod download && \
    go build -ldflags="-s -w" -o ./app ./cmd/hammer

# final stage
FROM gcr.io/distroless/base:nonroot
COPY --from=build-env /build/app /
COPY --from=build-env /build/db/migrations /db/migrations
EXPOSE 8000
EXPOSE 50051
ENTRYPOINT ["/app"]
