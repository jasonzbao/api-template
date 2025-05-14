-include .env
export

docker_build:
	docker buildx build --platform linux/arm64 -t api -f docker/api/Dockerfile .

docker_push:
	@bash ./scripts/docker_push.sh -e $(ENV) -a $(AWS_ACCOUNT_ID)

TIMEOUT ?= 300
deploy: docker_build docker_push
	DB_URI_FOR_ECS=""; \
	if [ "$(ENV)" = "prod" ]; then \
		DB_URI_FOR_ECS="$(PROD_DB_URI)"; \
	elif [ "$(ENV)" = "dev" ]; then \
		DB_URI_FOR_ECS="$(DEV_DB_URI)"; \
	else \
		echo "Error: Unsupported ENV value '$(ENV)' for DB_URI. Please set ENV to 'prod' or 'dev'." >&2; \
		exit 1; \
	fi; \
	ecs deploy $(ENV)-main $(ENV)-api \
		-e api VERSION $(shell git rev-parse --short HEAD) \
		-e api DB_URI "$${DB_URI_FOR_ECS}" \
		--timeout ${TIMEOUT} --user "$(shell id -un)"

migrate_generate:
	migrate create -ext sql -dir db/migrations -seq $(NAME)
