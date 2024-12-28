#!/bin/bash

set -e

DOCKER_TAG="latest"

TAG=$(git describe --tags --exact-match 2> /dev/null) || echo ""

if [[ -n "$TAG" ]]; then
    DOCKER_TAG="stable"
fi

echo "Building ${DOCKER_TAG} version"

docker build -t pagu-discord:${DOCKER_TAG}  -f ./deployment/Dockerfile . --target discord
docker build -t pagu-telegram:${DOCKER_TAG} -f ./deployment/Dockerfile . --target telegram

docker compose -f ./deployment/docker-compose.yml down
docker compose -f ./deployment/docker-compose.yml up -d

## Some cleanup
echo "Cleanup"

docker builder prune -f
docker image prune -f
