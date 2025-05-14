FROM golang:1.24 AS build

COPY ../go.mod /usr/src
COPY ../go.sum /usr/src

WORKDIR /usr/src

RUN go mod download

COPY ../../api /usr/src/api
COPY ../../Makefile /usr/src

RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags jsoniter -o cmd api/cmd/main.go

FROM golang:1.24 as migrate_builder

ARG MIGRATE_VERSION=v4.18.0
ENV CGO_ENABLED=0
ENV GOBIN=/go/bin
ENV PATH=${GOBIN}:${PATH}

# Install necessary tools and then golang-migrate
RUN <<EOF
set -ex # Exit on error, print commands
apt-get update
apt-get install -y --no-install-recommends ca-certificates git

echo "Installing migrate CLI version ${MIGRATE_VERSION} with postgres drivers..."
# Add -v for verbose output from go install
go install -v -tags 'postgres' "github.com/golang-migrate/migrate/v4/cmd/migrate@${MIGRATE_VERSION}"

echo "Migrate CLI installed to ${GOBIN}/migrate"
ls -l "${GOBIN}/migrate"

echo "Verifying migrate installation..."
"${GOBIN}/migrate" --version

# Cleanup
apt-get purge -y --auto-remove git
rm -rf /var/lib/apt/lists/*
EOF

FROM alpine:latest

# This is needed for the healthcheck in ECS
RUN apk --no-cache add curl
RUN apk add --no-cache tzdata

# Install tini
RUN apk add --no-cache tini

COPY --from=build /usr/src/cmd /usr
COPY --from=migrate_builder /go/bin/migrate /usr/local/bin/migrate
RUN migrate --version
COPY ../../scripts/docker_exec.sh /usr
COPY ../../configs /usr/configs
COPY ../../db/migrations /usr/db/migrations

WORKDIR /usr
RUN chmod +x docker_exec.sh

ENTRYPOINT ["/sbin/tini", "--", "/usr/docker_exec.sh"]
