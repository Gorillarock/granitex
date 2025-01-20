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

mocks:
	mockery

mock-clean:  ## Remove all generated mocks
	find "./" -type f -name "*_mock.go" -exec rm {} +

mockery-full-rebuild: mock-clean mocks
	go install github.com/vektra/mockery/v2@2.51.0
	@mockery --config=.mockery.yaml

test:
	go clean -testcache && go test -v ./...