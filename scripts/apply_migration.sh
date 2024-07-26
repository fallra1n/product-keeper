#!/usr/bin/env bash

CUR_DIR=$(dirname $0)

export DB_NAME=postgres
export DB_PASSWORD=pass

if [ -z "$1" ]
  then
    echo "Укажите номер миграции"
    exit 1
fi

migrate -path ${CUR_DIR?}/../migrations -database "postgres://${DB_NAME?}:${DB_PASSWORD?}@localhost:5432/${DB_NAME?}?sslmode=disable" up $1
