FROM debian:buster-slim AS build-env

# set workdir
WORKDIR /build

# install curl
RUN apt-get update && apt-get install -o Dpkg::Options::='--force-confdef' -o Dpkg::Options::='--force-confold' -y curl

# install golang-migrate
RUN curl -sfL https://github.com/golang-migrate/migrate/releases/download/v4.11.0/migrate.linux-amd64.tar.gz | tar -xvz && \
    echo "d0e186d92314c05fa842f9a9c32fd5fb14c6fc93c097766140fad563222aa325  migrate.linux-amd64" | sha256sum -c

# final stage
FROM gcr.io/distroless/base:nonroot
COPY --from=build-env /build/migrate.linux-amd64 /go-migrate
COPY db/migrations /tmp/
ENTRYPOINT ["/go-migrate"]
