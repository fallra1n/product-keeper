export NAME=product_keeper

SHELL := /bin/bash

mock:
	./scripts/automock.sh

no_test_cache:
	go clean -testcache

test_core: no_test_cache
	
test_adapter: no_test_cache
	
build:
	cd deployment && docker compose build --no-cache

run: build
	cd deployment && docker compose -p ${NAME} up --force-recreate --remove-orphans --build

.PHONY: build run
