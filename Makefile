export NAME=product_keeper

SHELL := /bin/bash
CORE_PATH := ./internal/core/...

mock:
	./scripts/automock.sh

no_test_cache:
	go clean -testcache

test_core: no_test_cache
	go test ${CORE_PATH} -v | grep -v "no test files"

test_adapter: no_test_cache 
	./scripts/run_test_adapter.sh

build:
	cd infra && docker compose build --no-cache --build-arg JWT_SECRET=${JWT_SECRET} --build-arg CONFIG_PATH=${CONFIG_PATH}

# Run only infrastructure (postgres and kafka)
run-infra:
	cd infra && docker compose -f docker-compose.infra.yaml -p ${NAME}-infra up --force-recreate --remove-orphans

# Run everything (original behavior)
run: test_core test_adapter build
	cd infra && docker compose -p ${NAME} up --force-recreate --remove-orphans

.PHONY: no_test_cache run-infra
