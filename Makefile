# Makefile

include .env

postgresinit:
	@docker-compose up -d

postgres:
	@docker exec -it $(POSTGRES_CONTAINER) psql -U $(DB_USER)

createdb:
	@docker exec -it $(POSTGRES_CONTAINER) psql -U $(DB_USER) -c "CREATE DATABASE $(DB_NAME);"

dropdb:
	@docker exec -it $(POSTGRES_CONTAINER) psql -U $(DB_USER) -c "DROP DATABASE IF EXISTS $(DB_NAME);"

migrate3partclear:
	@migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose force 3

migrateup:
	@migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up

migratedown:
	migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down

api_run:
	@docker run --rm --name $(CONTAINER_NAME) -p $(PORT):$(PORT) $(IMAGE_NAME)

api_build:
	@docker build -t $(IMAGE_NAME) .
	@echo "Docker image created: $(IMAGE_NAME)"

api_stop:
	@if [ "$(shell docker ps -q -f name=$(CONTAINER_NAME))" ]; then \
		docker stop $(CONTAINER_NAME); \
		docker rm $(CONTAINER_NAME); \
		echo "Docker container stopped: $(CONTAINER_NAME)"; \
	else \
		echo "No running container with name: $(CONTAINER_NAME)"; \
	fi

clean: stop
	@docker rmi $(IMAGE_NAME) || true
	@echo "Docker image dropped: $(IMAGE_NAME)"


.PHONY: postgresinit postgres createdb dropdb migrateup migratedown
