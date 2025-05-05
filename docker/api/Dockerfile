FROM golang:1.24 AS build

COPY ../go.mod /usr/src
COPY ../go.sum /usr/src

WORKDIR /usr/src

RUN go mod download

COPY ../../api /usr/src/api
COPY ../../Makefile /usr/src

RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags jsoniter -o cmd api/cmd/main.go

FROM alpine:latest

# This is needed for the healthcheck in ECS
RUN apk --no-cache add curl
RUN apk add --no-cache tzdata

# Install tini
RUN apk add --no-cache tini

COPY --from=build /usr/src/cmd /usr
COPY ../../scripts/docker_exec.sh /usr
COPY ../../api/configs /usr/configs

WORKDIR /usr
RUN chmod +x docker_exec.sh

ENTRYPOINT ["/sbin/tini", "--", "/usr/docker_exec.sh"]
