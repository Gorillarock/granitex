.PHONY: full-build

rebuild: mod clean force-build

force-build:
	docker compose build --no-cache

mod:
	go mod download

clean:
	go mod tidy

build: mod clean 
	docker compose build

run-with-rebuild: rebuild
	docker compose up

run: build
	docker compose up