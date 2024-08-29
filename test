#!/bin/bash

# Проверяем, был ли передан аргумент 'build'
if [ "$1" == "build" ]; then
  BUILD_OPTION="--build"
else
  BUILD_OPTION=""
fi

docker container stop db
docker container rm db
docker compose \
  -f compose.yaml \
  -f compose-test.yaml \
  up --abort-on-container-exit \
  --attach backend-testing $BUILD_OPTION


