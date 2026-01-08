#!/usr/bin/env bash

CUR_DIR=$(dirname $0)

if [ -z "$1" ]
  then
    echo "Укажите название миграции"
    exit 1
fi

NAME=$1

migrate create -ext sql -dir ${CUR_DIR?}/../migrations -seq $NAME
