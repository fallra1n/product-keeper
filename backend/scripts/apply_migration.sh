#!/usr/bin/env bash

CUR_DIR=$(dirname $0)

export DB_NAME=postgres
export DB_PASSWORD=pass
export DB_PORT=${DB_PORT:-5432}

if [ -z "$1" ]
then
  echo "enter the number of migration"
  exit 1
fi

MAX_RETRIES=5
RETRY_COUNT=0
SUCCESS=0

while [ ${RETRY_COUNT?} -lt ${MAX_RETRIES?} ]; do
  migrate -path ${CUR_DIR?}/../migrations -database "postgres://${DB_NAME?}:${DB_PASSWORD?}@localhost:${DB_PORT?}/${DB_NAME?}?sslmode=disable" up $1
  
  if [ $? -eq 0 ]; then
    SUCCESS=1
    break
  else
    sleep 20
    RETRY_COUNT=$((RETRY_COUNT+1))
  fi
done

if [ $SUCCESS -ne 1 ]; then
  echo "failed to apply migration"
  exit 1
fi
