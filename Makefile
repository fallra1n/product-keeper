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
	cd deployment && docker compose build --no-cache

run: test_core test_adapter build
	cd deployment && docker compose -p ${NAME} up --force-recreate --remove-orphans --build

.PHONY: no_test_cache
