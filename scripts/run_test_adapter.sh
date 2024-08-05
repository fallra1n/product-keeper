#!/usr/bin/env bash

ADAPTER_PATH=./internal/adapters/...

export CONFIG_PATH=../../../../config/local.yaml
export DB_PORT=5433

function setup_db() {
  echo "checking and removing existing postgres-test container if necessary..."
  if docker ps -a --filter "name=postgres-test" --format "{{.Names}}" | grep -q "^postgres-test$"; then
    docker stop postgres-test
    docker rm postgres-test
  fi
  docker compose up -d test-db
}

function remove_db() {
  echo "removing created db"
  if docker ps -a --filter "name=postgres-test" --format "{{.Names}}" | grep -q "^postgres-test$"; then
    docker stop postgres-test
    docker rm postgres-test
  fi
}

function apply_migrations() {
  echo "Applying migrations..."
  ./scripts/apply_migration.sh 1
}

cd deployment
setup_db
cd ..

apply_migrations

echo
echo ==================================================
echo

go test ${ADAPTER_PATH} | grep -v "no test files"

echo
echo ==================================================

remove_db