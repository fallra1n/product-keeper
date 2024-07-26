export NAME=product_keeper

SHELL := /bin/bash

build:
	cd deployment && docker compose build --no-cache

run: build
	cd deployment && docker compose -p ${NAME} up --force-recreate --remove-orphans --build

.PHONY: build run
