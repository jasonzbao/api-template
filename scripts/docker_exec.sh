#!/bin/sh

/usr/local/bin/migrate -database ${DB_URI} -path db/migrations up

exec ./cmd --config ${CONFIG_FILE} --version ${VERSION}
