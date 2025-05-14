-include .env
export

docker_build:
	docker buildx build --platform linux/arm64 -t api -f docker/api/Dockerfile .

docker_push:
	@bash ./scripts/docker_push.sh -e $(ENV) -a $(AWS_ACCOUNT_ID)

TIMEOUT ?= 300
SKIP_MIGRATE ?= 0
deploy: docker_build docker_push
	if [ $(SKIP_MIGRATE) = 0 ]; then\
		make migrate ENV=$(ENV);\
	fi
	ecs deploy $(ENV)-gcc $(ENV)-api -e api VERSION $(shell git rev-parse --short HEAD) --timeout ${TIMEOUT} --user "$(shell id -un)"

migrate:
	ifeq ($(ENV), prod)
		migrate -database "$(PROD_DB_URI)" -path db/migrations up
	else ifeq ($(ENV), dev)
		migrate -database "$(DEV_DB_URI)" -path db/migrations up
	else
		@echo "Error: Unsupported ENV value '$(ENV)'. Please set ENV to 'prod' or 'dev'."
		@exit 1
	endif

migrate_generate:
	migrate create -ext sql -dir db/migrations -seq $(NAME)
