docker_build:
	docker buildx build --platform linux/arm64 -t api -f docker/Dockerfile .

docker_push:
	@bash ./scripts/docker_push.sh -e $(ENV) -a $(AWS_ACCOUNT_ID)

TIMEOUT ?= 300
deploy: docker_build docker_push
	if [ $(SKIP_MIGRATE) = 0 ]; then\
		make migrate ENV=$(ENV);\
	fi
	ecs deploy $(ENV)-gcc $(ENV)-api -e api VERSION $(shell git rev-parse --short HEAD) --timeout ${TIMEOUT} --user "$(shell id -F)"

migrate:
	migrate -database $(DB_URI) -path db/migrations up

migrate_generate:
	migrate create -ext sql -dir db/migrations -seq $(NAME)
