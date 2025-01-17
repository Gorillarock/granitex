.PHONY: build

run: build
	docker compose up

run-headless: build
	docker compose up -d

run-with-rebuild: rebuild
	docker compose up

run-headless-with-rebuild: rebuild
	docker compose up -d

build:
	docker compose build

rebuild:
	docker compose build --no-cache
